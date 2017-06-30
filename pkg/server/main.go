package server

import (
	"github.com/jiusanzhou/pdf2html/pkg/server/backend"
	"github.com/jiusanzhou/tentacle/log"
	"os/exec"
	"os"
	"syscall"
	"time"
	"github.com/jiusanzhou/pdf2html/pkg/util"
)

// Type of service signals currently supported.
type serviceSignal int

const (
	serviceStatus  = iota // Gets status about the service.
	serviceRestart        // Restarts the service.
	serviceStop           // Stops the server.
	// Add new service requests here.
)

const (
	serverShutdownPoll = 500 * time.Millisecond
)

// Global service signal channel.
var globalServiceSignalCh chan serviceSignal

// Global service done channel.
var globalServiceDoneCh chan struct{}

// Initialize service mutex once.
func init() {
	globalServiceDoneCh = make(chan struct{}, 1)
	globalServiceSignalCh = make(chan serviceSignal)
}

// restartProcess starts a new process passing it the active fd's. It
// doesn't fork, but starts a new process using the same environment and
// arguments as when it was originally started. This allows for a newly
// deployed binary to be started. It returns the pid of the newly started
// process when successful.
func restartProcess() error {
	// Use the original binary location. This works with symlinks such that if
	// the file it points to has been changed we will use the updated symlink.
	argv0, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}

	// Pass on the environment and replace the old count key with the new one.
	cmd := exec.Command(argv0, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

// Handles all serviceSignal and execute service functions.
func (server *Server) handleServiceSignals() error {

	// Wait for SIGTERM in a go-routine.
	trapCh := util.SignalTrap(os.Interrupt, syscall.SIGTERM)
	go func(trapCh <-chan bool) {
		<-trapCh
		globalServiceSignalCh <- serviceStop
	}(trapCh)

	// Start listening on service signal. Monitor signals.
	for {
		signal := <-globalServiceSignalCh
		switch signal {
		case serviceStatus:
			/// We don't do anything for this.
		case serviceRestart:
			if err := server.stop(); err != nil {
				server.Error("Unable to close server gracefully.", err)
			}
			if err := restartProcess(); err != nil {
				server.Error("Unable to restart the server.", err)
			}
		case serviceStop:
			server.Info("Received signal to exit.")
			if err := server.stop(); err != nil {
				server.Error("Unable to close server gracefully.", err)
			}
		}
	}
}

type Server struct {

	backends []backend.Backend

	currentJobs int32

	config *Config

	log.Logger

	shutdown *util.Shutdown

	// forcibly close them during graceful stop or restart.
	gracefulTimeout time.Duration
}

func (server *Server) backend()backend.Backend {
	return server.backends[0]
}

func (server *Server) run() {

}

func (server *Server) close() {

	server.shutdown.WaitBegin()

	// clean every thing

	server.shutdown.Complete()
}

func (server *Server) stop() error {

	server.Info("Wait for [%d]jobs to complete.", server.currentJobs)

	// send msg to shutdown for begin
	// who is wait shutdown begin
	// should go to the next
	server.shutdown.Begin()

	// some body should call Complete
	// process go-routin should use
	// wait begin and complete to
	// handle shutdown

	server.shutdown.WaitComplete()

	server.Info("Exit server success.")

	return nil

}

func (server *Server) Start() {

}

func NewServer(c Config)(server *Server, err error) {

	server = &Server{}

	server.gracefulTimeout = 60 * time.Minute

	return
}