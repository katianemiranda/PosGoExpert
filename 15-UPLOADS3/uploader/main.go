package main

import (
	"fmt"
	"os"
)

var (
	s3Client *s3.S3
	s3Bucket string
	wg       sync.WaitGroup
)

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func init() {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String("us-west-1"),
			Credentials: credentials.NewStaticCredentials(
				"your-access-key-id", 
				"your-secret", 
				""
			),
		},
	)
	if err != nil {
		panic(err)
	}
	s3Client = s3.New(sess)
	s3Bucket = "goexpert-bucket-exemplo"
}

func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	uploadControl := make(chan struct{}, 100)
	errorFileUpload := make(chan string)

	go func(){
		for {
			select {
			case fileName := <-errorFileUpload:
				wg.Add(1)
				uploadControl <- struct{}{}
				uploadFile(fileName, uploadControl, errorFileUpload)
			}
		}
	}

	for {
		files, err := dir.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading directory: %v\n", err)
			continue
		}
		wg.Add(1)
		uploadControl <- struct{}{}
		uploadFile(files[0].Name(), uploadControl)
	}

}


func uploadFile(fileName string, uploadControl <-chan struct{}, errorFileUpload chan<- string) {
	defer wg.Done()
	completeFileName := fmt.Sprintf("./tmp/%s", fileName)
	fmt.Print("Uploading file: ", completeFileName, "\n")
	f, err := os.Open(completeFileName)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", completeFileName, err)
		<-uploadControl
		errorFileUpload <- fileName
		return
	}
	defer f.Close()

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fileName),
		Body:   f,
	})
	if err != nil {
		fmt.Printf("Error uploading file %s: %v\n", completeFileName, err)
		<-uploadControl
		errorFileUpload <- fileName
		return
	}
	fmt.Printf("Successfully uploaded %s to %s\n", completeFileName, s3Bucket)
	<-uploadControl
}