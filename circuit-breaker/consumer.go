package main

import (
	"context"
	"time"

	"github.com/sony/gobreaker/v2"
)

type Consumer struct {
	breaker *gobreaker.CircuitBreaker[any]
	repo    *Repository
}

func NewConsumer(breaker *gobreaker.CircuitBreaker[any], repo *Repository) *Consumer {
	return &Consumer{
		breaker: breaker,
		repo:    repo,
	}
}

func (c *Consumer) Process(ctx context.Context, order Order) error {
	_, err := c.breaker.Execute(func() (any, error) {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		return nil, c.repo.Save(ctx, order)
	})

	return err
}
