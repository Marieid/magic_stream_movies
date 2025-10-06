package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	//This is the main function

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, Magic_stream_movies!")
	})
	
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Filed to start server", err)
	}
}
