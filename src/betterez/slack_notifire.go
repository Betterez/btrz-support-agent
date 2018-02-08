package betterez

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	slackWebhookURL = fmt.Sprintf("https://hooks.slack.com/services/%s/%s/%s",
		os.Getenv("SLACK_TEAM"),
		os.Getenv("SLACK_GROUP"),
		os.Getenv("SLACK_TOKEN"))
)

const (
	contentType = "application/json"
	username    = "SupportDB1"
	channel     = "#project"
)

// SendSlackNotification send a message to the notification channel
func SendSlackNotification(message string) bool {
	fmt.Print("using this utl:", slackWebhookURL)
	dataString := fmt.Sprintf(`{"channel":"%s","icon_emoji":":satellite_antenna:", "username": "%s", "text": "%s"}`,
		channel,
		username,
		message)
	fmt.Print(dataString)
	resp, err := http.Post(slackWebhookURL, contentType, strings.NewReader(dataString))
	if resp.StatusCode < 400 && nil == err {
		return true
	}
	return false
}
