package main

import (
	"net/http"
	route "presigned_url/modules"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := route.SetupRoutes()
	port := ":8080"

	if err := router.Run(port); err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running on " + port))
	})
}
