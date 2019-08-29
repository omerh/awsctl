package hooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type slackRequestBody struct {
	Text string `json:"text"`
}

func getWebhookURLFromEnvironmentVariable() string {
	if value, ok := os.LookupEnv("AWSCTL_SLACK_URL"); ok {
		return value
	}
	fmt.Println("Missing environment variable AWSCTL_SLACK_URL")
	return "none"
}

// SendSlackWebhook for sending custom message to slack for monitoring
func SendSlackWebhook(message string) {
	webhookURL := getWebhookURLFromEnvironmentVariable()
	if webhookURL == "none" {
		return
	}
	slackBody, _ := json.Marshal(slackRequestBody{Text: message})

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		fmt.Println(errors.New("Non-ok response returned from Slack"))
	}
}
