package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sony/gobreaker/v2"
)

type Worker struct {
	breaker *gobreaker.CircuitBreaker[any]
	repo    *Repository
}

func NewWorker(breaker *gobreaker.CircuitBreaker[any], repo *Repository) *Worker {
	return &Worker{
		breaker: breaker,
		repo:    repo,
	}
}

func (c *Worker) Process(ctx context.Context, msg string) error {
	var order Order
	if err := json.Unmarshal([]byte(msg), &order); err != nil {
		return fmt.Errorf("Failed to unmarshal SQS message body: %w", err)
	}

	_, err := c.breaker.Execute(func() (any, error) {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		return nil, c.repo.Save(ctx, order)
	})

	return err
}
