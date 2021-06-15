package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/betterez/btrz-support-agent/utils"
	"github.com/bsphere/le_go"
)

const (
	CurrentVersion = "1.0.0.9"
)

func pullerService() {
	leToken := os.Getenv("LE_TOKEN")
	deploymentData := utils.CreateDeploymentFromEnvVars()
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
	le.Close()
	for {
		le, err = le_go.Connect(leToken)
		if err != nil {
			log.Fatalf("can't use log entries inside for loop. Error %v", err)
		} else {
			response, err := utils.GetSQSMessage(awsSession)
			if err != nil {
				le.Print(err, " while trying to pull sqs message ", CurrentVersion)
			}
			if response != nil {
				le.Print(CurrentVersion, "got new support request")
				le.Close()
				utils.SendSlackNotification("New db requests received, follow LE for execution details.")
				deploymentData.DatabaseName = response.DatabaseName
				ok, lastObject, err := utils.LoadLastObjectFromBucket(response.BucketName, "dump.tar.gz", awsSession)
				le, _ = le_go.Connect(leToken)
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
				le.Close()
				ok, err = utils.DeployToServer("dump.tar.gz", deploymentData)
				le, _ = le_go.Connect(leToken)
				if err != nil {
					le.Printf("Error %v while deploying object, bucket %s", err, response.BucketName)
					continue
				}
				if !ok {
					le.Printf("failed deploying object , bucket %s (no error)", response.BucketName)
					continue
				}
				le.Print("deployed, anonymizing")
				le.Close()
				ok, err = utils.AnonymizeDB(deploymentData)
				le, _ = le_go.Connect(leToken)
				if err != nil {
					le.Printf("Error %v while anonymizing db %s", err, response.DatabaseName)
					continue
				}
				if !ok {
					le.Printf("failed anonymizing db , db %s (no error)", response.DatabaseName)
					continue
				}
				le.Print("done anonymizing.")
				utils.SendSlackNotification("Operation completed!")
			}
		}
		le.Close()
		time.Sleep(time.Minute)
	}
}
func showLastObjectSelected() {
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		log.Fatalf("error %v", err)
	}
	object, err := utils.GetLastObject(os.Getenv("BUCKET_NAME"), awsSession)
	if err != nil {
		log.Fatal("error ", err)
	}
	fmt.Print(*object.Key)
}
func main() {
	pullerService()
}
