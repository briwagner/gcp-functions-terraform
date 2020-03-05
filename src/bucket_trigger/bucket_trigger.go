package bucket_trigger

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
)

type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// Global API client.
var storageClient *storage.Client

// Environment variables
var topicID string
var projectID string

func BucketTrigger(ctx context.Context, e GCSEvent) error {
	topicID = os.Getenv("TOPIC_ID")
	projectID = os.Getenv("PROJECT_ID")

	// Process new file upload.
	uri := fmt.Sprintf("gs://%s/%s", e.Bucket, e.Name)
	msg := fmt.Sprintf("Received a new file %s", uri)
	log.Print(msg)

	// Process file as needed...

	// Send message on pubsub topic.
	if topicID != "" {
		err := SendMessage(ctx, msg, "success", topicID, projectID)
		if err != nil {
			log.Printf("Unable to send pubsub message %v", err)
			return err
		}
	}
	return nil
}

// SendMessage sends the message on the given pubsub topic.
func SendMessage(ctx context.Context, m string, status string, topic string, project string) error {
	log.Printf("Sending message on topic %s", topic)
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return err
	}

	// Add custom attributes.
	attrs := map[string]string{
		"status": status,
	}
	msg := &pubsub.Message{
		Data:       []byte(m),
		Attributes: attrs,
	}

	// Publish message to topic.
	t := client.Topic(topic)
	if _, err := t.Publish(ctx, msg).Get(ctx); err != nil {
		return err
	}

	return nil
}
