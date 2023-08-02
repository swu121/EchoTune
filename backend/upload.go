package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Generate a random userID
func generateUserID() string {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", b)
}

func getS3url(c *gin.Context) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	))

	// Create a new instance of the service's client with a Session.
	svc := s3.New(sess)

	userID := generateUserID()

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("echolab"),
		Key:    aws.String(userID + "/"),
	})

	if err != nil {
		log.Fatalf("Unable to create folder %q, %v", userID, err)
	}

	fmt.Printf("Successfully created folder %q\n", userID)

	// Now we generate the presigned URL for the specific file within the 'folder'
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("echolab"),
		Key:    aws.String(userID + "/warmup"),
	})
	urlStr, err := req.Presign(15 * 60)

	if err != nil {
		log.Println("Failed to sign request", err)
	}

	log.Println("Presigned URL is", urlStr)

	c.IndentedJSON(http.StatusOK, urlStr)
}

// func upload(url string) {
// 	uploader := s3manager.NewUploader(sess)

// 	f, err := os.Open("../warmupVocal1.wav")
// 	if err != nil {
// 		log.Fatalf("failed to open file %q, %v", "../warmupVocal1.wav", err)
// 	}

// 	// Upload the file to S3.
// 	_, err = uploader.Upload(&s3manager.UploadInput{
// 		Bucket: aws.String("echolab"),
// 		Key:    aws.String(userID + "/warmup"),
// 		Body:   f,
// 	})

// 	if err != nil {
// 		log.Fatalf("failed to upload file, %v", err)
// 	}

// 	log.Println("file uploaded successfully")
// }
