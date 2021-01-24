package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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

	router.POST("/api/file", func(ctx *gin.Context) {
		f, uploadedFile, err := ctx.Request.FormFile("file")
		if err != nil {
			ErrorHandler(err, ctx)
			return
		}
		defer f.Close()

		signedURL, token, err := uploadToGCS(ctx, &f, uploadedFile.Filename)
		if err != nil {
			ErrorHandler(err, ctx)
			return
		}

		err = SetURLAndToken(ctx, token, signedURL)
		if err != nil {
			ErrorHandler(err, ctx)
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"url": os.Getenv("ENTRYPOINT") + "/" + token,
		})
	})

	router.GET("/:token", func(ctx *gin.Context) {
		token := ctx.Param("token")
		signedURL, err := GetURLFromToken(ctx, token)
		if err != nil {
			ErrorHandler(err, ctx)
			return
		}
		ctx.Redirect(http.StatusFound, signedURL)
	})

	//router.NoRoute(func(ctx *gin.Context) {
	//	ctx.File("client/")
	//})

	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("cannot start gin: %v", err)
	}
}

func ErrorHandler(err error, ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"message": err.Error(),
		"error":   true,
	})
}
