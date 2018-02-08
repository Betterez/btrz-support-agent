package betterez

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	sqsURL      = "https://sqs.us-east-1.amazonaws.com/109387325558/btrz-support-agent"
	sqsPullType = "All"
)

// RequestInformation - sqs request from json data
type RequestInformation struct {
	BucketName   string `json:"bucket_name"`
	DatabaseName string `json:"db_name"`
}

// GetSQSMessage gets the pull messages from the queue
func GetSQSMessage(session *session.Session) (*RequestInformation, error) {
	sqsService := sqs.New(session)
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(sqsURL), // Required
		AttributeNames: []*string{
			aws.String(sqsPullType),
		},
		MaxNumberOfMessages: aws.Int64(1),
		MessageAttributeNames: []*string{
			aws.String(sqsPullType),
		},
		VisibilityTimeout: aws.Int64(1),
		WaitTimeSeconds:   aws.Int64(1),
	}
	resp, err := sqsService.ReceiveMessage(params)
	if err != nil {
		return nil, err
	}
	if len(resp.Messages) == 0 {
		return nil, nil
	}
	deleteParam := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(sqsURL),
		ReceiptHandle: resp.Messages[0].ReceiptHandle,
	}
	sqsService.DeleteMessage(deleteParam)
	requestData := *resp.Messages[0].Body
	jsonDecoder := json.NewDecoder(strings.NewReader(requestData))
	var result = &RequestInformation{}
	jsonDecoder.Decode(result)
	return result, nil
}
