package signal

import (
	"context"
	"os"
	ossignal "os/signal"
	"syscall"
)

var defaultQuitSignals = []os.Signal{os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM}

func ListenQuit(parent context.Context) context.Context {
	return Listen(parent, defaultQuitSignals...)
}

func Listen(parent context.Context, signals ...os.Signal) context.Context {
	ch := make(chan os.Signal, 1)
	ossignal.Notify(ch, signals...)

	ctx, cancel := context.WithCancel(parent)

	go func() {
		defer cancel()
		defer ossignal.Stop(ch)

		select {
		case <-parent.Done():
		case <-ch:
		}
	}()

	return ctx
}
