package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/upload", getS3url)
	router.Run("localhost:8080")
}
