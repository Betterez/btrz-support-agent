package utils

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetLastObject getting the last object in a bucket
func GetLastObject(bucketName string, sess *session.Session) (*s3.Object, error) {
	var selectedObject *s3.Object
	var batchSize int64 = 500
	svc := s3.New(sess)
	listingParams := &s3.ListObjectsInput{
		Bucket:    aws.String(bucketName), // Required
		Delimiter: aws.String("Delimiter"),
		MaxKeys:   &batchSize,
	}
	for {
		resp, err := svc.ListObjects(listingParams)
		if err != nil {
			return nil, err
		}
		for _, currentObject := range resp.Contents {
			if selectedObject == nil {
				selectedObject = currentObject
			} else {
				if currentObject.LastModified.After(*selectedObject.LastModified) {
					selectedObject = currentObject
				}
			}
		}
		if len(resp.Contents) < int(batchSize) {
			break
		} else {
			listingParams.Marker = selectedObject.Key
		}
	}
	return selectedObject, nil
}

func loadBackup(selectedObject *s3.Object, bucketName, destinationFileName string, sess *session.Session) error {
	svc := s3.New(sess)
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    selectedObject.Key,
	}
	resp, err := svc.GetObject(params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	destinationFile, err := os.Create(destinationFileName)
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	dataBytes := make([]byte, 4096)
	for {
		bytesRead, _ := resp.Body.Read(dataBytes)
		if err != nil {
			return err
		}
		if bytesRead > 0 {
			bytesToWrite := dataBytes[:bytesRead]
			destinationFile.Write(bytesToWrite)
		} else {
			break
		}
	}
	return nil
}

//LoadLastObjectFromBucket loading last object from a bucket in to destination file
func LoadLastObjectFromBucket(bucketName, destinationFileName string, sess *session.Session) (bool, *s3.Object, error) {
	lastObject, err := GetLastObject(bucketName, sess)
	if err != nil {
		return false, nil, err
	}
	if lastObject == nil {
		return false, nil, nil
	}
	err = loadBackup(lastObject, bucketName, destinationFileName, sess)
	if err != nil {
		return false, nil, err
	}
	return true, lastObject, nil
}
