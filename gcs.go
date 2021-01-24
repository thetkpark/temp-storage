package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"
	"time"
)

func uploadToGCS(f *multipart.File, fileName string) error {
	var bucketName = os.Getenv("BUCKET_NAME")
	token, err := generateUniqueToken()
	if err != nil {
		return fmt.Errorf("generateUniqueToken: %v", err)
	}
	var prefix = strconv.FormatInt(time.Now().Unix(), 10) + "-" + token
	objectFullpath := prefix + "/" + fileName

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("gcs-sa-key.json"))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucketName).Object(objectFullpath).NewWriter(ctx)
	if _, err = io.Copy(wc, *f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil

	// TODO: Get signed URL

	// TODO: Return signedURL and token
}

func getSignedURL() (string, error) {
	var bucketName = os.Getenv("BUCKET_NAME")
	object := "notes.txt"

	pkey, err := ioutil.ReadFile("gcs-sa-key.pem")
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	url, err := storage.SignedURL(bucketName, object, &storage.SignedURLOptions{
		GoogleAccessID: os.Getenv("GOOGLE_ACCESS_ID"),
		PrivateKey:     pkey,
		Method:         "GET",
		Expires:        time.Now().Add(72 * time.Hour),
	})
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}
	return url, nil
}
