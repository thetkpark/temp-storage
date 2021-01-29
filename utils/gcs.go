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
	"strconv"
	"time"
)

func UploadToGCS(ginContext context.Context, f *bytes.Buffer, token string, fileName string) (string, error) {
	var bucketName = os.Getenv("BUCKET_NAME")

	var prefix = strconv.FormatInt(time.Now().Unix(), 10) + "-" + token
	objectFilepath := prefix + "/" + fileName

	client, err := storage.NewClient(ginContext, option.WithCredentialsFile("gcs-sa-key.json"))
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ginContext, time.Second*30)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucketName).Object(objectFilepath).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	// TODO: Get signed URL
	signedUrl, err := getSignedURL(objectFilepath)
	if err != nil {
		return "", fmt.Errorf("getSignedURL: %v", err)
	}

	// TODO: Return signedURL and token
	return signedUrl, nil
}

func getSignedURL(object string) (string, error) {
	var bucketName = os.Getenv("BUCKET_NAME")

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
