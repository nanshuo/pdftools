package util

import (
	"os"
	"os/signal"
)

// signalTrap traps the registered signals and notifies the caller.
func SignalTrap(sig ...os.Signal) <-chan bool {
	// channel to notify the caller.
	trapCh := make(chan bool, 1)

	go func(chan<- bool) {
		// channel to receive signals.
		sigCh := make(chan os.Signal, 1)
		defer close(sigCh)

		// `signal.Notify` registers the given channel to
		// receive notifications of the specified signals.
		signal.Notify(sigCh, sig...)

		// Wait for the signal.
		<-sigCh

		// Once signal has been received stop signal Notify handler.
		signal.Stop(sigCh)

		// Notify the caller.
		trapCh <- true
	}(trapCh)

	return trapCh
}
