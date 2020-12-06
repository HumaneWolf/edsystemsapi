package api

import (
	"github.com/gin-gonic/gin"
	"humanewolf.com/ed/systemapi/systems"
)

// RunAPI starts the program's rest API.
func RunAPI() {
	defer systems.CloseFiles()
	r := gin.Default()

	// Default includes these already.
	//r.Use(gin.Logger())
	//r.Use(gin.Recovery())
	r.Use(CheckAccessControl())

	r.GET("/typeahead", handleTypeahead)
	r.GET("/index_stats", handleIndexStats)

	go systems.StartCacheCleaner()
	r.Run()
}
