package betterez

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bitly/go-simplejson"
)

var (
	configFileName = "./config/notifications.json"
	sqsPullURL     = ""
	sqsPushURL     = ""
	sqsPullType    = "All"
	serverKeyFile  = ""
)

// RequestInformation - sqs request from json data
type RequestInformation struct {
	BucketName   string `json:"bucket_name"`
	DatabaseName string `json:"db_name"`
}

func init() {
	log.Println("loading config data")
	if os.Getenv("CONFIG_FILE") != "" {
		configFileName = os.Getenv("CONFIG_FILE")
	}
	if err := loadConfigData(configFileName); err != nil {
		log.Println("loading data error:", err)
	} else {
		log.Println("done!")
	}
}

func loadConfigData(fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	data, err := simplejson.NewFromReader(file)
	if err != nil {
		return err
	}
	defer file.Close()
	sqsPullURL, err = data.Get("incoming").String()
	if err != nil {
		sqsPullURL = ""
		return err
	}
	log.Printf("sqsPullURL=%s\r\n", sqsPullURL)
	sqsPushURL, _ = data.Get("outgoing").String()
	log.Printf("sqsPushURL=%s\r\n", sqsPushURL)
	serverKeyFile, _ = data.Get("server-key").String()
	log.Printf("serverKeyFile=%s\r\n", serverKeyFile)
	if serverKeyFile != "" {
		if _, err = os.Stat(serverKeyFile); os.IsNotExist(err) {
			log.Printf("key file %s does not exist", serverKeyFile)
		} else {
			log.Printf("key file %s found!", serverKeyFile)
		}
	}
	return nil
}

// GetSQSMessage gets the pull messages from the queue
func GetSQSMessage(session *session.Session) (*RequestInformation, error) {
	if sqsPullURL == "" {
		return nil, errors.New("No input string for sqs message")
	}
	sqsService := sqs.New(session)
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(sqsPullURL), // Required
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
		QueueUrl:      aws.String(sqsPullURL),
		ReceiptHandle: resp.Messages[0].ReceiptHandle,
	}
	sqsService.DeleteMessage(deleteParam)
	requestData := *resp.Messages[0].Body
	jsonDecoder := json.NewDecoder(strings.NewReader(requestData))
	var result = &RequestInformation{}
	jsonDecoder.Decode(result)
	return result, nil
}

// SetCompletionMessage - send completion message
func SetCompletionMessage(awsSession *session.Session, messageData string) (*sqs.SendMessageOutput, error) {
	if sqsPushURL == "" {
		return nil, errors.New("no output queue")
	}
	sqsService := sqs.New(awsSession)
	response, err := sqsService.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageBody:  aws.String(messageData),
		QueueUrl:     aws.String(sqsPushURL),
	})
	return response, err
}
