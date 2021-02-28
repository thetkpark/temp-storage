package controllers

import (
	"encoding/base64"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thetkpark/tempStorage/utils"
)

func UploadFileController(ctx *gin.Context) {
	f, uploadFile, err := ctx.Request.FormFile("file")
	if err != nil {
		errorHandler(err, ctx)
		return
	}
	if uploadFile.Size > 524288000 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "File size larget than 500MB",
		})
		return
	}
	defer f.Close()

	if uploadFile.Size > 104857600 {
		ctx.JSON(400, gin.H{
			"message": "File is larger than 100MB",
			"error":   true,
		})
		return
	}
	// Generate encryption key
	key, err := utils.GenerateEncryptionKey()
	if err != nil {
		errorHandler(err, ctx)
	}

	// Generate retrieve token key
	token, err := utils.GenerateUniqueToken()
	if err != nil {
		errorHandler(err, ctx)
	}

	ObjectName, err := utils.GenerateFileName()
	if err != nil {
		errorHandler(err, ctx)
	}

	// Create buffer
	encryptedBuffer, err := utils.EncryptFile(&f, key)
	if err != nil {
		errorHandler(err, ctx)
	}

	// Upload to GCS and get signedURL
	err = utils.UploadToGCS(ctx, encryptedBuffer, ObjectName)
	if err != nil {
		errorHandler(err, ctx)
		return
	}

	// Set URL and token in Redis
	var fileData = utils.FileMetadata{
		FileName:   base64.StdEncoding.EncodeToString([]byte(uploadFile.Filename)),
		Key:        key,
		ObjectName: ObjectName,
	}
	err = utils.SetTokenFileData(ctx, token, fileData)
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

	fileData, err := utils.GetFileDataFromToken(ctx, token)
	if err != nil {
		if err.Error() == "rdb.Get: redis: nil" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Not found",
				"error":   true,
			})
			return
		}
		errorHandler(err, ctx)
		return
	}
	fileName, err := base64.StdEncoding.DecodeString(fileData.FileName)
	if err != nil {
		errorHandler(err, ctx)
		return
	}

	encryptedFile, err := utils.DownloadFile(ctx, fileData.ObjectName)
	if err != nil {
		errorHandler(err, ctx)
		return
	}

	decryptedFile := utils.DecryptFile(encryptedFile, fileData.Key)

	ctx.Header("Content-Disposition", `attachment; filename="`+string(fileName)+`"`)
	ctx.Data(200, "application/octet-stream", *decryptedFile)
}

func errorHandler(err error, ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"message": err.Error(),
		"error":   true,
	})
}
