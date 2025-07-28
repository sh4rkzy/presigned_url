package route

import (
	"presigned_url/modules/presigned/controller"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "AWS S3 Presigned URL Service",
		})
	})

	r.GET("/s3/presigned-url", controller.GeneratePresignedURL)
	r.GET("/s3/presigned-get", controller.GeneratePresignedGET)
	
	return r
}
