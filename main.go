package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thetkpark/tempStorage/controllers"
	"log"
	"os"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "5000"
	}

	// Setup Gin
	router := gin.Default()

	router.POST("/api/file", controllers.UploadFileController)

	router.GET("/:token", controllers.GetFileController)

	//router.NoRoute(func(ctx *gin.Context) {
	//	ctx.File("client/")
	//})

	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("cannot start gin: %v", err)
	}
}


