package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// load .env file from given path
	// we keep it empty it will load .env from current directory
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	fmt.Println()
	//err := uploadToGCS()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//url, err := getSignedURL()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(url)
}