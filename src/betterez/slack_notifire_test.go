package betterez

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type SlackInfo struct {
	Team  string `json:"slack_team"`
	Group string `json:"slack_group"`
	Token string `json:"slack_token"`
}

func TestChannelPosting(t *testing.T) {
	t.SkipNow()
	testFileName := "../mocks/slack_data.json"
	if _, err := os.Stat(testFileName); os.IsNotExist(err) {
		t.SkipNow()
	}
	slackData, err := os.Open(testFileName)
	if err != nil {
		t.Fatal("can't process file data", err)
	}
	jsonDecoder := json.NewDecoder(slackData)
	slackInfo := &SlackInfo{}
	jsonDecoder.Decode(slackInfo)
	slackWebhookURL = fmt.Sprintf("https://hooks.slack.com/services/%s/%s/%s",
		slackInfo.Team,
		slackInfo.Group,
		slackInfo.Token)
	flag := SendSlackNotification("test")
	if flag == false {
		t.Fatal("Slack did not get this message!")
	}
}
