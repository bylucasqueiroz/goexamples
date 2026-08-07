package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sony/gobreaker/v2"

	pgx "github.com/jackc/pgx/v5"
)

const (
	baseDelay  = 1 * time.Second
	maxRetries = 5
)

func main() {
	ctx := context.Background()
	conn, err := dbConnect()
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
	consumer := NewConsumer(breaker, repo)

	worker(ctx, conn, queue, consumer)
}

func worker(ctx context.Context, conn *pgx.Conn, queue *Queue, consumer *Consumer) error {
	queueUrl := os.Getenv("SQS_QUEUE_URL")

	for {
		messages, err := queue.GetMessage(ctx, queueUrl)
		if err != nil {
			return fmt.Errorf("Error receiving messages: %v", err)
		}

		for _, msg := range messages {
			if msg.Body == nil {
				log.Println("Message body is empty")
				continue
			}

			var order Order
			err := json.Unmarshal([]byte(*msg.Body), &order)
			if err != nil {
				log.Fatalf("Failed to unmarshal SQS message body: %v", err)
			}

			err = consumer.Process(ctx, order)
			switch {
			case err == nil:
				queue.DeleteMessage(ctx, queueUrl, *msg.ReceiptHandle)
			case errors.Is(err, gobreaker.ErrOpenState):
				log.Println("Circuit Breaker aberto")
			case errors.Is(err, gobreaker.ErrTooManyRequests):
				log.Println("Half-Open limit")
			default:
				log.Printf("erro ao salvar: %v", err)
			}
		}
	}
}

func dbConnect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v", err)
	}
	return conn, nil
}
