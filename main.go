package main

import (
	"example.com/lambda/entities"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"sync"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func Handler() (entities.Response, error) {

	sess, err := session.NewSession()

	if err != nil {
		exitErrorf(err.Error())
	}

	svc := s3.New(sess)
	bucket := os.Getenv("BUCKET_NAME")

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: &bucket,
	})

	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	contents := entities.CreateNewResponse()

	contents = updateBucketItems(resp, svc, bucket, contents)

	return contents, nil
}

func getString(item string) *string {
	return &item
}

func updateBucketItems(
	resp *s3.ListObjectsV2Output,
	svc *s3.S3,
	bucket string,
	contents entities.Response,
) entities.Response {
	var wg sync.WaitGroup
	wg.Add(len(resp.Contents))
	for i, item := range resp.Contents {
		item := item
		go func(i int) {
			defer wg.Done()

			_, err := svc.PutObject(&s3.PutObjectInput{
				Bucket:   &bucket,
				Key:      item.Key,
				Metadata: map[string]*string{"content": getString("ok")},
			})

			if err != nil {
				fmt.Printf("Error while trying to update object %q, %v", *item.Key, err)
				return
			}

			contents.AddItem(entities.S3Object{
				Name:         *item.Key,
				LastModified: *item.LastModified,
				Size:         *item.Size,
			})

			fmt.Printf("Finished processing %q. Going to another item on the list.", *item.Key)
		}(i)
	}

	wg.Wait()

	return contents
}

func main() {
	lambda.Start(Handler)
}
