package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thetkpark/tempStorage/utils"
)

func UploadFileController(ctx *gin.Context) {
	f, uploadedFile, err := ctx.Request.FormFile("file")
	if err != nil {
		errorHandler(err, ctx)
		return
	}
	defer f.Close()

	if uploadedFile.Size > 104857600 {
		ctx.JSON(400, gin.H{
			"message": "File is larger than 100MB",
			"error":   true,
		})
		return
	}
	// Generate retrieve token key
	token, err := utils.GenerateUniqueToken()
	if err != nil {
		errorHandler(err, ctx)
	}

	encryptedBuffer, err := utils.EncryptFile(&f)
	if err != nil {
		errorHandler(err, ctx)
	}

	// Upload to GCS and get signedURL
	signedURL, err := utils.UploadToGCS(ctx, encryptedBuffer, token, uploadedFile.Filename)
	if err != nil {
		errorHandler(err, ctx)
		return
	}

	// Set URL and token in Redis
	err = utils.SetURLAndToken(ctx, token, signedURL)
	if err != nil {
		errorHandler(err, ctx)
		return
	}

	// Return the retrieve token key
	ctx.JSON(http.StatusCreated, gin.H{
		"url": os.Getenv("ENTRYPOINT") + "/" + token,
	})
}

func GetFileController(ctx *gin.Context) {
	token := ctx.Param("token")
	signedURL, err := utils.GetURLFromToken(ctx, token)
	if err != nil {
		errorHandler(err, ctx)
		return
	}
	ctx.Redirect(http.StatusFound, signedURL)
}

func errorHandler(err error, ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"message": err.Error(),
		"error":   true,
	})
}
