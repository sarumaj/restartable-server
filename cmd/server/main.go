package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sarumaj/restartable-server/pkg/response"
	"github.com/sarumaj/restartable-server/pkg/restartableserver"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// record the timestamp, the server is being created at
	startedAt := time.Now()

	// create new server
	server := restartableserver.NewRestartableServer(
		&http.Server{
			Addr:      ":" + restartableserver.Getenv("PORT", "8080"),
			TLSConfig: nil,
		},
		log.Default(),
	)

	serverID := fmt.Sprintf("%p", server.GetServer())

	// create fallback server to handle incoming requests while productive server is down
	fallback := restartableserver.NewFallbackServer(server.GetServer())

	// dispatch fallback server until productive server is ready
	go func() { _ = fallback.ListenAndServe() }()

	// create and configure traffic handler
	router := gin.New()
	router.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			Output:    server.GetLogger().Writer(),
			SkipPaths: []string{"/uptime"},
		}),
		gin.Recovery(),
	)
	response.LoadHTMLTemplates(router)

	// configure routes
	{
		// serve "ok.html" at root
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "ok.html", gin.H{"ServerID": serverID})
		})
		// return uptime at "/uptime"
		router.GET("/uptime", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"Elapsed": restartableserver.FormatDuration(time.Since(startedAt)),
			})
		})
		// force server to restart at "/restart"
		router.GET("/restart", func(c *gin.Context) {
			defer func() {
				server.InterruptSender() <- os.Interrupt
			}()

			c.Status(http.StatusAccepted)
		})
	}

	// install handler on the server
	server.SetHandler(router)

	// terminate fallback server to release port binding for the productive server
	_ = fallback.Close()

	restartAfter := 10 * time.Second // time interval before the come back of the productive server
	shutdownTimeout := time.Minute   // graceful shutdown period

	// dispatch recovery from panics
	defer server.RecoverAndRestart(main, restartAfter)

	// run server
	go server.ListenAndServeWithRecover(main, restartAfter)

	// invoke callback for shutdown or restart
	server.ShutdownOrRestart(main, restartAfter, shutdownTimeout)
}
