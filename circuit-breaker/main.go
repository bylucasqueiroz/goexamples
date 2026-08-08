package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sony/gobreaker/v2"

	pgx "github.com/jackc/pgx/v5"
)

const (
	maxRetries = 5
	baseDelay  = 1 * time.Second
	maxDelay   = 10 * time.Second
)

func main() {
	ctx := context.Background()
	conn, err := initDbConnection()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	sdkConfig, err := config.LoadDefaultConfig(
		ctx, config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatalf("Error loading default configuration: %v\n", err)
	}
	client := sqs.NewFromConfig(sdkConfig, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("LOCALSTACK_ENDPOINT"))
	})

	queue := NewQueue(client)
	repo := NewRepository(conn)
	breaker := NewDatabaseBreaker()
	worker := NewWorker(breaker, repo)

	if err := consumer(ctx, queue, worker); err != nil {
		log.Fatalf("Worker failed: %v", err)
	}
}

func consumer(ctx context.Context, queue *Queue, worker *Worker) error {
	queueUrl := os.Getenv("SQS_QUEUE_URL")
	retryCount := 0

	for {
		messages, err := queue.GetMessage(ctx, queueUrl)
		if err != nil {
			return fmt.Errorf("Error receiving messages: %w", err)
		}

		if len(messages) == 0 {
			continue
		}

		shouldBackoff := false
		for _, msg := range messages {
			if msg.Body == nil {
				log.Println("Message body is empty")
				continue
			}

			err = worker.Process(ctx, *msg.Body)
			if err == nil {
				if err := queue.DeleteMessage(ctx, queueUrl, *msg.ReceiptHandle); err != nil {
					log.Printf("Failed to delete message from queue: %v", err)
				}
				log.Printf("Successfully processed order: %v", msg.Body)
				retryCount = 0
				continue
			}

			if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
				log.Println("Circuit Breaker aberto")
				shouldBackoff = true
				break
			}

			log.Printf("erro ao salvar: %v", err)
			break
		}

		if shouldBackoff {
			retryCount = doBackoff(ctx, retryCount)
		}
	}
}

func doBackoff(ctx context.Context, retryCount int) int {
	backoff := baseDelay * (1 << retryCount)
	delay := time.Duration(rand.Int64N(int64(backoff)))
	if delay > maxDelay {
		delay = maxDelay
	}
	log.Printf("Retrying in backoff %v...\n", delay)

	select {
	case <-time.After(delay):
	case <-ctx.Done():
	}

	if retryCount < maxRetries {
		retryCount++
	}
	return retryCount
}

func initDbConnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v", err)
	}
	return conn, nil
}
