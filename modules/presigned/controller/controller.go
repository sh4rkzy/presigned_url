package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	config_aws "presigned_url/infrastructure/aws"
	uuid "presigned_url/shared"
	"time"
)

type TransactionData struct {
	TransactionID string `json:"transactionID"`
	Timestamp     string `json:"timestamp"`
}

type PresignedResponse struct {
	Sucesso     int             `json:"sucesso"`
	Message     string          `json:"message"`
	URL         string          `json:"url"`
	Transaction TransactionData `json:"transaction"`
}

func GeneratePresignedURL(c *gin.Context) {
	bucket := c.Query("bucket")
	key := c.Query("key")

	if bucket == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   1,
			"message": "bucket and key parameters are required",
			"transaction": TransactionData{
				TransactionID: uuid.GenerateUUID(),
				Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
			},
		})
		return
	}

	url, err := config_aws.GeneratePresignedURL(config_aws.Config{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate presigned URL",
		})
		return
	}

	response := PresignedResponse{
		Sucesso: 0,
		Message: "Presigned URL generated successfully",
		URL:     url,
		Transaction: TransactionData{
			TransactionID: uuid.GenerateUUID(),
			Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	c.JSON(http.StatusOK, response)
}

func GeneratePresignedGET(c *gin.Context) {
	bucket := c.Query("bucket")
	key := c.Query("key")

	if bucket == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   1,
			"message": "bucket and key parameters are required",
			"transaction": TransactionData{
				TransactionID: uuid.GenerateUUID(),
				Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
			},
		})
		return
	}

	url, err := config_aws.GeneratePresignedGET(config_aws.Config{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate presigned GET URL",
		})
		return
	}

	response := PresignedResponse{
		Sucesso: 0,
		Message: "Presigned GET URL generated successfully",
		URL:     url,
		Transaction: TransactionData{
			TransactionID: uuid.GenerateUUID(),
			Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	c.JSON(http.StatusOK, response)
}
