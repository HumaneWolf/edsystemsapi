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

	// Default includes these already.
	//r.Use(gin.Logger())
	//r.Use(gin.Recovery())

	r.GET("/typeahead", func(c *gin.Context) {
		input, exists := c.GetQuery("input")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No input."})
			return
		}
		results := systems.SearchTreeForNames(strings.TrimSpace(input))
		c.JSON(200, results)
	})

	r.GET("/index_stats", func(c *gin.Context) {
		c.JSON(200, systems.GetIndexStats())
	})

	r.Run()
}
