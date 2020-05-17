package api

import (
	"github.com/gin-gonic/gin"
	"humanewolf.com/ed/systemapi/systems"
)

// RunAPI starts the program's rest API.
func RunAPI() {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/search", func(c *gin.Context) {
		input := c.Param("input")
		results := systems.SearchTree(input)
		c.JSON(200, results)
	})
	r.Run()
}
