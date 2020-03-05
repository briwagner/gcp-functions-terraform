package slackbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
)

var slackURL string

// SlackMessage is the outgoing message to Slack.
type SlackMessage struct {
	Text string `json:"text"`
}

// SendMessage is the function entrypoint.
func SendMessage(ctx context.Context, m pubsub.Message) error {
	// Set value from environment variable.
	slackURL = os.Getenv("SLACK_URL")

	// Exit if environment var not set.
	if slackURL == "" {
		log.Fatalf("Slack URL must be set as environment variable: slackURL.")
	}

	// Build message, using 'status' attribute to generate emoji.
	var emoji string

	if _, ok := m.Attributes["status"]; ok {
		switch m.Attributes["status"] {
		case "success":
			emoji = ":+1:"
		case "error":
			emoji = ":octagonal_sign:"
		default:
			emoji = ":heavy_check_mark:"
		}
	} else {
		// Default if no status is passed.
		emoji = ":heavy_check_mark:"
	}

	// Format message string.
	message := SlackMessage{
		Text: fmt.Sprintf("%s %s", emoji, string(m.Data)),
	}

	// Build request to Slack.
	b, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Could not encode json: %v", err)
		return err
	}
	req, err := http.NewRequest("POST", slackURL, bytes.NewBuffer(b))
	if err != nil {
		log.Fatalf("Could not build post request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Could not make post request: %v", err)
		return err
	}

	// Log response from Slack.
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		// Error response is plain text.
		rr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Could not parse response from Slack: %s", err)
			return err
		}
		log.Printf("Slack message failed: %v", string(rr))
		return nil
	}

	log.Print("Slack message completed.")
	return nil
}
