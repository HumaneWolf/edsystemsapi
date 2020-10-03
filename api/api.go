package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"humanewolf.com/ed/systemapi/systems"
)

// RunAPI starts the program's rest API.
func RunAPI() {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/search", func(c *gin.Context) {
		input, exists := c.GetQuery("input")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No input."})
			return
		}
		results := systems.SearchTree(strings.TrimSpace(input))
		c.JSON(200, results)
	})
	r.Run()
}
