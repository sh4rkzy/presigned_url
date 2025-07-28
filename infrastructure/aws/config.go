package config_aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Bucket string
	Key    string
	Region string
}

func GeneratePresignedURL(cfg Config) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	region := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),

		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			os.Getenv("AWS_SESSION_TOKEN"),
		),
	})
	if err != nil {
		log.Println("Failed to create session:", err)
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(cfg.Key),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Failed to sign request:", err)
		return "", err
	}

	return urlStr, nil
}

func GeneratePresignedGET(cfg Config) (string, error) {
	region := os.Getenv("AWS_REGION")
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println("Failed to create session:", err)
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(cfg.Key),
	})
	return req.Presign(15 * time.Minute)
}
