package api

import "github.com/gin-gonic/gin"

// RunAPI starts the program's rest API.
func RunAPI() {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
