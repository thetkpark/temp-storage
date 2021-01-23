package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"os"
	"time"
)

var bucketName = os.Getenv("BUCKET_NAME")


func uploadToGCS() error {

	token, err := generateUniqueToken()
	if err != nil {
		return fmt.Errorf("generateUniqueToken: %v", err)
	}
	var prefix = string(time.Now().Unix()) + "-" + token
	object := "notes.txt"
	objectFullpath := prefix + object

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("gcs-sa-key.json"))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open("notes.txt")
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucketName).Object(objectFullpath).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Println("Uploaded")
	return nil
}

func getSignedURL () (string, error) {
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