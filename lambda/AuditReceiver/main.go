package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"
)

const DELIVERYSTREAMNAME = "kdf-firehose-audit-folgado"

type Actor struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
}

type Action struct {
	Resource  string `json:"resource"`
	Operation string `json:"operation"`
	Status    string `json:"status"`
}

type Context struct {
	IPAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent"`
}

type Change struct {
	Field    string `json:"field"`
	OldValue string `json:"oldValue"`
	NewValue string `json:"newValue"`
}

type Details struct {
	EntityType string   `json:"entityType"`
	EntityID   string   `json:"entityId"`
	ChangeType string   `json:"changeType"`
	Changes    []Change `json:"changes"`
}

type Metadata struct {
	ProcessedAt time.Time `json:"processedAt"`
	Source      string    `json:"source"`
}

type Event struct {
	ID         string   `json:"_id"`
	EventType  string   `json:"eventType"`
	EntityType string   `json:"entityType"`
	EntityID   string   `json:"entityId"`
	Timestamp  string   `json:"timestamp"`
	RequestID  string   `json:"requestId"`
	Actor      Actor    `json:"actor"`
	Action     Action   `json:"action"`
	Context    Context  `json:"context"`
	Details    Details  `json:"details"`
	Metadata   Metadata `json:"metadata"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event Event
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}

	eventResult, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to process event",
		}, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Failed to load AWS config: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to load AWS configuration",
		}, nil
	}

	firehoseClient := firehose.NewFromConfig(cfg)

	firehosePayload := firehose.PutRecordInput{
		DeliveryStreamName: aws.String(DELIVERYSTREAMNAME),
		Record: &types.Record{
			Data: eventResult,
		},
	}

	_, err = firehoseClient.PutRecord(context.TODO(), &firehosePayload)

	if err != nil {
		log.Printf("Failed to put record in Firehose: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to put record in Firehose",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Event processed successfully",
	}, nil
}

func main() {
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return Handler(request)
	})
}
