package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	controller "github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/controllers"
)

func main() {
	//This is the main function

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, Magic_stream_movies!")
	})

	router.GET("/movies", controller.GETMovies())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Filed to start server", err)
	}
}
