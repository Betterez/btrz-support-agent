package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/betterez/btrz-support-agent/utils"
)

const (
	CurrentVersion = "1.0.11"
)

func pullerService() {
	deploymentData := utils.CreateDeploymentFromEnvVars()
	if !deploymentData.IsLegal() {
		log.Fatalf("Can't get mongo variables!")
	}
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		log.Fatalf("error %v", err)
	}
	log.Print("Support server running version ", CurrentVersion)
	//le.Close()
	for {
		response, err := utils.GetSQSMessage(awsSession)
		if err != nil {
			log.Print(err, " while trying to pull sqs message ", CurrentVersion)
		}
		if response != nil {
			log.Print(CurrentVersion, "got new support request")

			utils.SendSlackNotification("New db requests received, follow LE for execution details.")
			deploymentData.DatabaseName = response.DatabaseName
			ok, lastObject, err := utils.LoadLastObjectFromBucket(response.BucketName, "dump.tar.gz", awsSession)

			if err != nil {
				log.Printf("Error %v while loading object from s3, bucket %s", err, response.BucketName)
				continue
			}
			if !ok {
				log.Printf("failed loading object from s3, bucket %s (no error)", response.BucketName)
				continue
			}
			log.Print(*lastObject.Key, " was loaded")
			log.Print("object loaded from s3, deploying")

			ok, err = utils.DeployToServer("dump.tar.gz", deploymentData)
			if err != nil {
				log.Printf("Error %v while deploying object, bucket %s", err, response.BucketName)
				continue
			}
			if !ok {
				log.Printf("failed deploying object , bucket %s (no error)", response.BucketName)
				continue
			}
			log.Print("deployed, anonymizing")

			ok, err = utils.AnonymizeDB(deploymentData)
			if err != nil {
				log.Printf("Error %v while anonymizing db %s", err, response.DatabaseName)
				continue
			}
			if !ok {
				log.Printf("failed anonymizing db , db %s (no error)", response.DatabaseName)
				continue
			}
			log.Print("done anonymizing.")
			utils.SendSlackNotification("Operation completed!")
		}
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
