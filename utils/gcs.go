package utils

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func UploadToGCS(ginContext context.Context, f *bytes.Buffer, fileName string)  error {
	var bucketName = os.Getenv("BUCKET_NAME")

	client, err := storage.NewClient(ginContext, option.WithCredentialsFile("gcs-sa-key.json"))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ginContext, time.Second*30)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	// TODO: Return signedURL and token
	return nil
}

//func getSignedURL(object string) (string, error) {
//	var bucketName = os.Getenv("BUCKET_NAME")
//
//	pkey, err := ioutil.ReadFile("gcs-sa-key.pem")
//	if err != nil {
//		return "", fmt.Errorf("ioutil.ReadFile: %v", err)
//	}
//	url, err := storage.SignedURL(bucketName, object, &storage.SignedURLOptions{
//		GoogleAccessID: os.Getenv("GOOGLE_ACCESS_ID"),
//		PrivateKey:     pkey,
//		Method:         "GET",
//		Expires:        time.Now().Add(72 * time.Hour),
//	})
//	if err != nil {
//		return "", fmt.Errorf("storage.SignedURL: %v", err)
//	}
//	return url, nil
//}

func DownloadFile(ctx context.Context,object string) (*[]byte, error) {
	var bucketName = os.Getenv("BUCKET_NAME")

	client, err := storage.NewClient(ctx, option.WithCredentialsFile("gcs-sa-key.json"))
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucketName).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	return &data, nil
}