package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

var quitSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT}

func QuitSignal() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, quitSignals...)

	go func() {
		defer close(ch)
		defer signal.Stop(ch)

		<-ch
	}()

	return ch
}
