package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thetkpark/tempStorage/controllers"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")

	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "5000"
	}

	clientPath := os.Getenv("CLIENT_PATH")
	if len(clientPath) < 1 {
		clientPath = "client/"
	}

	// Setup Gin
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST"},
	}))

	// router.StaticFile("/index.html", "/client/index.html")

	router.POST("/api/file", controllers.UploadFileController)

	router.GET("/:token", controllers.GetFileController)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.File(clientPath)
	})

	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("cannot start gin: %v", err)
	}
}
