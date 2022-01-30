package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
)

//AWS_ACCESS_KEY_ID
//AWS_SECRET_ACCESS_KEY
var maxCopies = 0

func main() {

	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Fatalln("environment variable `AWS_SECRET_ACCESS_KEY` not set")
	}

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		log.Fatalln("environment variable `AWS_ACCESS_KEY_ID` not set")
	}

	if os.Getenv("AWS_REGION") == "" {
		log.Fatalln("environment variable `AWS_REGION` not set")
	}

	if os.Getenv("AWS_ENDPOINT") == "" {
		log.Fatalln("environment variable `AWS_ENDPOINT` not set")
	}

	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("AWS_ENDPOINT")

	if len(os.Args) < 4 {
		log.Fatalln(`usage: $> s3backup bucket directory filetobackup [maxCopies]`)
	}

	bucket := os.Args[1]
	directory := os.Args[2]
	filename := os.Args[3]

	if len(os.Args) >= 5 {
		var err error
		maxCopies, err = strconv.Atoi(os.Args[4])
		if err != nil {
			log.Fatalln("invalid max copies value, should be integer")
		}
	}

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	svc := s3.New(sess)
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(directory),
	}

	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("start uploading file", path.Base(filename))
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String("private"),
		Key:    aws.String(path.Join(directory, path.Base(filename))),
		Body:   file,
	})

	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("uploaded file", up.Location)

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Fatalln(err.Error())
	}

	//remove the oldest files in directory
	if maxCopies > 0 && len(resp.Contents) > maxCopies {
		log.Printf("removing oldest backups as there are only %d copies allowed", maxCopies)
		//sort them recent date to the least recent
		sort.Slice(resp.Contents, func(i, j int) bool {
			return resp.Contents[i].LastModified.After(*resp.Contents[j].LastModified)
		})

		for _, key := range resp.Contents[maxCopies:] {
			log.Printf("Deleting old backup `%s` \n", *key.Key)
			_, err := svc.DeleteObject(&s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    key.Key,
			})
			if err != nil {
				log.Printf("Error deleting old backup `%s` with error `%s`\n", *key.Key, err.Error())
			}
		}
	}

	log.Println("s3backup is done")
}
