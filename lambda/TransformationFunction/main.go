package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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
	ProcessedAt string `json:"processedAt"`
	Source      string `json:"source"`
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

func Handler(ctx context.Context, firehoseEvent events.KinesisFirehoseEvent) (events.KinesisFirehoseResponse, error) {
	var transformedRecords []events.KinesisFirehoseResponseRecord

	for _, record := range firehoseEvent.Records {
		var event Event

		err := json.Unmarshal(record.Data, &event)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			transformedRecords = append(transformedRecords, events.KinesisFirehoseResponseRecord{
				RecordID: record.RecordID,
				Result:   events.KinesisFirehoseTransformedStateDropped,
				Data:     []byte{},
			})
			continue
		}

		eventJSON, err := json.Marshal(event)
		if err != nil {
			log.Printf("Error marshaling event to JSON: %v", err)
			transformedRecords = append(transformedRecords, events.KinesisFirehoseResponseRecord{
				RecordID: record.RecordID,
				Result:   events.KinesisFirehoseTransformedStateDropped,
				Data:     []byte{},
			})
			continue
		}

		transformedData := append(eventJSON, '\n')

		var eventTime time.Time
		if event.Timestamp != "" {
			parsedTime, err := time.Parse("2006-01-02:15:04:05Z", event.Timestamp)
			if err != nil {
				log.Printf("Error parsing timestamp, using current time: %v", err)
				eventTime = time.Now().UTC()
			} else {
				eventTime = parsedTime
			}
		} else {
			eventTime = time.Now().UTC()
		}

		year := eventTime.Format("2006")
		month := eventTime.Format("01")
		day := eventTime.Format("02")

		transformedRecords = append(transformedRecords, events.KinesisFirehoseResponseRecord{
			RecordID: record.RecordID,
			Result:   events.KinesisFirehoseTransformedStateOk,
			Data:     transformedData,
			Metadata: events.KinesisFirehoseResponseRecordMetadata{
				PartitionKeys: map[string]string{
					"audit_type": event.EventType,
					"year":       year,
					"month":      month,
					"day":        day,
				},
			},
		})
	}

	return events.KinesisFirehoseResponse{Records: transformedRecords}, nil
}

func main() {
	lambda.Start(Handler)
}
