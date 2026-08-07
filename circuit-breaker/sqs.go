package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Queue struct {
	SqsClient *sqs.Client
}

func NewQueue(sqsClient *sqs.Client) *Queue {
	return &Queue{
		SqsClient: sqsClient,
	}
}

func (a *Queue) GetMessage(ctx context.Context, queueUrl string) ([]types.Message, error) {
	var messages []types.Message
	result, err := a.SqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		log.Printf("Couldn't get messages from queue %v. Here's why: %v\n", queueUrl, err)
	} else {
		messages = result.Messages
	}
	return messages, err
}

func (a *Queue) DeleteMessage(ctx context.Context, queueUrl string, receiptHandle string) error {
	_, err := a.SqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		log.Printf("Couldn't delete message from queue %v. Here's why: %v\n", queueUrl, err)
	}
	return err
}
