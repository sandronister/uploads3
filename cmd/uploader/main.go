package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sandronister/uploads3/configs"
)

var (
	s3Client   *s3.S3
	s3Bucket   string
	fileErrors []string
)

func init() {
	awsConfig, err := configs.GetConfig(".")

	if err != nil {
		panic(err)
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			awsConfig.Credentials,
			awsConfig.Password,
			"",
		),
	})

	if err != nil {
		panic(err)
	}

	s3Client = s3.New(sess)
	s3Bucket = awsConfig.Bucket
}

func sendFile(file <-chan string) {
	for filename := range file {
		completeFileName := fmt.Sprintf("./tmp/%s", filename)
		f, err := os.Open(completeFileName)
		defer f.Close()

		if err != nil {
			fileErrors = append(fileErrors, filename)
			fmt.Printf("Erro ao abrir arquivo %s\n", filename)
			continue
		}

		fmt.Printf("Arquivo aberto %s\n", filename)

		_, err := s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(filename),
			Body:   f,
		})

		if err != nil {
			fileErrors = append(fileErrors, filename)
		}
	}
}

func getFiles(fileChan chan<- string) {
	files, err := os.ReadDir("./tmp")
	if err != nil {
		panic(err)
	}
	for _, item := range files {
		fileChan <- item.Name()
	}
}

func main() {

	limit := 100
	fileChan := make(chan string, 10)

	for i := 0; i < limit; i++ {
		go sendFile(fileChan)
	}

	getFiles(fileChan)
}
