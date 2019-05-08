package betterez

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"testing"
)

func TestSQSOutput(t *testing.T) {
	t.Skip()
	configFileName = "../../config/notifications.json"
	loadConfigData(configFileName)
	messageToSend := `{"type":"completion","service":"connex","environment":"support"}`
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		t.Fatalf("error %v", err)
	}
	_, err = SetCompletionMessage(awsSession, messageToSend)
	if err != nil {
		t.Fatalf("error %v", err)
	}
}
