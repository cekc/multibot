package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var DefaultQuitSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT}

func Context() context.Context {
	return Listen(context.Background(), DefaultQuitSignals...)
}

func Listen(parent context.Context, signals ...os.Signal) context.Context {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)

	ctx, cancel := context.WithCancel(parent)

	go func() {
		defer cancel()
		defer signal.Stop(ch)

		select {
		case <-parent.Done():
		case <-ch:
		}
	}()

	return ctx
}
