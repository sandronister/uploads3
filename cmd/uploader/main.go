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

	fmt.Print(awsConfig)

	if err != nil {
		panic(err)
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			awsConfig.AcessKey,
			awsConfig.SecretKey,
			"",
		),
	})

	if err != nil {
		panic(err)
	}

	s3Client = s3.New(sess)
	s3Bucket = awsConfig.Bucket
}

func uploads3(filename, shortname string) error {
	f, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer f.Close()
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(shortname),
		Body:   f,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Subiu arquivo %s\n", shortname)
	return nil
}

func sendFile(file <-chan string) {
	for filename := range file {
		completeFileName := fmt.Sprintf("../../tmp/%s", filename)

		err := uploads3(completeFileName, filename)

		if err != nil {
			fileErrors = append(fileErrors, filename)
			fmt.Printf("Erro ao abrir arquivo %s\n", filename)
		}

	}
}

func getFiles(fileChan chan<- string) {
	files, err := os.ReadDir("../../tmp")
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
