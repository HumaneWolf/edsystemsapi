package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"humanewolf.com/ed/systemapi/systems"
)

func handleTypeahead(c *gin.Context) {
	input, exists := c.GetQuery("input")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No input."})
		return
	}
	results := systems.SearchTreeForNames(strings.TrimSpace(input))
	c.JSON(200, results)
}

func handleIndexStats(c *gin.Context) {
	c.JSON(200, systems.GetIndexStats())
}
