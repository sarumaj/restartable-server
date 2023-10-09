package restartableserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Restartable server
type RestartableServer interface {
	InterruptReceiver() <-chan os.Signal
	InterruptSender() chan<- os.Signal
	GetLogger() *log.Logger
	GetServer() *http.Server
	ListenAndServeWithRecover(main func(), after time.Duration)
	RecoverAndRestart(main func(), after time.Duration)
	Restart(main func(), after time.Duration)
	SetHandler(handler http.Handler)
	ShutdownOrRestart(main func(), after, timeout time.Duration)
}

// should implement RestartableServer
type restartableServer struct {
	*http.Server
	*log.Logger
	sig chan os.Signal
}

// make sure restartableServer implements RestartableServer
var _ RestartableServer = new(restartableServer)

// Get underlying logger
func (s restartableServer) GetLogger() *log.Logger { return s.Logger }

// Get underlying HTTP server
func (s restartableServer) GetServer() *http.Server { return s.Server }

// Dispatches restartable server and handles interrupt signals and panics by restarting
func (s restartableServer) ListenAndServeWithRecover(main func(), after time.Duration) {
	defer s.RecoverAndRestart(main, after)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Panicln(err)
	}
}

// Exposes receiving only signal channel
func (s restartableServer) InterruptReceiver() <-chan os.Signal { return s.sig }

// Exposes sending only signal channel
func (s restartableServer) InterruptSender() chan<- os.Signal { return s.sig }

// Recover from panic by restarting the server
func (s restartableServer) RecoverAndRestart(main func(), after time.Duration) {
	if v := recover(); v != nil {
		s.Println(v)
		s.Restart(main, after)
	}
}

// Restarts the server by calling the main routine after specified amount of time
func (s restartableServer) Restart(main func(), after time.Duration) {
	later := time.Now().Add(after) // record the latter time point for timer
	tickFrequency := after / 10    // frequency at which the ticks will be logged

	if tickFrequency > time.Second*10 || tickFrequency == 0 {
		tickFrequency = time.Second * 10
	}

	// temporary fallback to handle count-down
	tmpFallback := NewFallbackServer(s.Server)

	// dispatch fallback server
	go func() { _ = tmpFallback.ListenAndServe() }()

	signal.Notify(s.InterruptSender(), os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(tickFrequency)
	timer := time.NewTimer(after)
	for {
		select {

		case t := <-ticker.C: // sending ticks to logger
			s.Printf("Restarting in %s... \n", later.Sub(t).Round(tickFrequency))

		case <-s.InterruptReceiver(): // terminate fallback server on interrupt
			ticker.Stop()
			_ = tmpFallback.Close()
			return

		case <-timer.C: // restart server by calling the main routine
			ticker.Stop()
			_ = tmpFallback.Close()

			defer main()
			return

		}
	}
}

// Set handler engine for the HTTP server
func (s restartableServer) SetHandler(handler http.Handler) { s.Server.Handler = handler }

// Gracefully shutdown the server with a graceful timeout
// and restart it after specified amount of time
func (s restartableServer) ShutdownOrRestart(main func(), after, timeout time.Duration) {
	defer s.Restart(main, after)

	signal.Notify(s.InterruptSender(), os.Interrupt, syscall.SIGTERM)
	s.Printf("Received signal: %v\n", <-s.InterruptReceiver())
	s.Println("Terminating server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		s.Panicln(err)
	}
}

// Create new restartable server
func NewRestartableServer(server *http.Server, logger *log.Logger) *restartableServer {
	s := &restartableServer{
		Server: server,
		Logger: logger,
		sig:    make(chan os.Signal, 1),
	}

	// make sure, an instance of logger is always available
	if s.Logger == nil {
		s.Logger = log.Default()
	}

	// default HTTP server
	if s.Server == nil {
		s.Server = &http.Server{}
	}

	return s
}
