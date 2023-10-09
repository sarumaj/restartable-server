package restartableserver

import (
	"fmt"
	"io"
	"net/http"
	"time"

	gin "github.com/gin-gonic/gin"
	response "github.com/sarumaj/restartable-server/pkg/response"
)

// Example implementation of a fallback server
func NewFallbackServer(s *http.Server) *http.Server {
	timestamp := time.Now()

	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: io.Discard}), gin.Recovery())
	response.LoadHTMLTemplates(router)

	server := &http.Server{
		Addr:      s.Addr,
		TLSConfig: s.TLSConfig,
		Handler:   router,
	}
	serverID := fmt.Sprintf("%p", server)

	{
		router.GET("/downtime", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"Elapsed": FormatDuration(time.Since(timestamp)),
			})
		})

		router.NoRoute(func(c *gin.Context) {
			c.HTML(http.StatusGone, "gone.html", gin.H{"ServerID": serverID})
		})
	}

	return server
}
