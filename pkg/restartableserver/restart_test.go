package restartableserver

import (
	"bytes"
	"log"
	"os"
	"testing"
	"time"

	gin "github.com/gin-gonic/gin"
)

func TestRestartableServer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := log.Default()
	buffer := bytes.NewBuffer(nil)
	logger.SetOutput(buffer)

	var server *restartableServer
	var main func()

	main = func() {
		server = NewRestartableServer(nil, logger)
		fallback := NewFallbackServer(server.GetServer())
		go func() { _ = fallback.ListenAndServe() }()

		router := gin.Default()
		server.SetHandler(router)

		_ = fallback.Close()
		go server.ListenAndServeWithRecover(main, time.Second)
		server.ShutdownOrRestart(main, time.Second, time.Minute)
	}

	go func() {
		<-time.After(time.Millisecond * 200)
		server.InterruptSender() <- os.Interrupt
		<-time.After(time.Millisecond * 1200)
		server.InterruptSender() <- os.Interrupt
		<-time.After(time.Millisecond * 200)
		server.InterruptSender() <- os.Interrupt
		<-time.After(time.Millisecond * 200)
		server.InterruptSender() <- os.Interrupt
	}()

	main()

	t.Log(buffer.String())
}
