package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"time"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

type S3Object struct {
	Name         string    `json:"name"`
	LastModified time.Time `json:"last_modified"`
	Size         int64     `json:"size"`
}

type Response struct {
	Items []S3Object `json:"items"`
}

func Handler() (Response, error) {

	sess, err := session.NewSession()

	if err != nil {
		exitErrorf(err.Error())
	}

	svc := s3.New(sess)
	bucket := os.Getenv("BUCKET_NAME")

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: &bucket})

	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	contents := Response{}

	// transform it into async
	for _, item := range resp.Contents {
		_, err := svc.PutObject(&s3.PutObjectInput{
			Bucket:  &bucket,
			Key:      item.Key,
			Metadata: map[string]*string{"foo" : aws.String("ok")},
		})

		if err != nil {
			log.Printf("Error while trying to update object %q, %v", *item.Key, err)
			continue
		}

		contents.Items = append(contents.Items, S3Object{
			Name:         *item.Key,
			LastModified: *item.LastModified,
			Size:         *item.Size,
		})
	}

	return contents, nil
}

func main() {
	lambda.Start(Handler)
}
