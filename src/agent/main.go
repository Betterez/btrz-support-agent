package main

import (
	"betterez"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bsphere/le_go"
)

const (
	CurrentVersion = "1.0.0.7"
)

func pullerService() {
	leToken := os.Getenv("LE_TOKEN")
	deploymentData := betterez.CreateDeploymentFromEnvVars()
	if !deploymentData.IsLegal() {
		log.Fatalf("Can't get mongo variables!")
	}
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		log.Fatalf("error %v", err)
	}
	le, err := le_go.Connect(leToken)
	if err != nil {
		log.Fatalf("can't use log entries. Error %v", err)
	}
	le.Print("Support server running version ", CurrentVersion)
	for {
		response, err := betterez.GetSQSMessage(awsSession)
		if err != nil {
			le.Print(err, " while trying to pull sqs message ", CurrentVersion)
		}
		if response != nil {
			le.Close()
			le, _ = le_go.Connect(leToken)
			le.Print(CurrentVersion, "got new support request")
			betterez.SendSlackNotification("New db requests received, follow LE for execution details.")
			deploymentData.DatabaseName = response.DatabaseName
			ok, lastObject, err := betterez.LoadLastObjectFromBucket(response.BucketName, "dump.tar.gz", awsSession)
			if err != nil {
				le.Printf("Error %v while loading object from s3, bucket %s", err, response.BucketName)
				continue
			}
			if !ok {
				le.Printf("failed loading object from s3, bucket %s (no error)", response.BucketName)
				continue
			}
			le.Print(*lastObject.Key, " was loaded")
			le.Print("object loaded from s3, deploying")
			ok, err = betterez.DeployToServer("dump.tar.gz", deploymentData)
			if err != nil {
				le.Printf("Error %v while deploying object, bucket %s", err, response.BucketName)
				continue
			}
			if !ok {
				le.Printf("failed deploying object , bucket %s (no error)", response.BucketName)
				continue
			}
			le, _ = le_go.Connect(leToken)
			le.Print("deployd, anonymizing")
			ok, err = betterez.AnonymizeDB(deploymentData)
			if err != nil {
				le.Printf("Error %v while anonymizing db %s", err, response.DatabaseName)
				continue
			}
			if !ok {
				le.Printf("failed anonymizing db , db %s (no error)", response.DatabaseName)
				continue
			}
			le, _ = le_go.Connect(leToken)
			le.Print("done anonymizing.")
			betterez.SendSlackNotification("Operation completed!")
		}
		time.Sleep(time.Minute)
	}
}
func showLastObjectSelected() {
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		log.Fatalf("error %v", err)
	}
	object, err := betterez.GetLastObject(os.Getenv("BUCKET_NAME"), awsSession)
	if err != nil {
		log.Fatal("error ", err)
	}
	fmt.Print(*object.Key)
}
func main() {
	pullerService()
}
