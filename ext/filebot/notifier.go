package filebot

import (
	"bufio"
	"context"
	"io"
)

type Notifier struct {
	writer bufio.Writer
	lock   chan struct{}
}

func CreateNotifier(writer io.Writer) Notifier {
	return Notifier{
		writer: *bufio.NewWriter(writer),
		lock:   make(chan struct{}, 1),
	}
}

func (notifier *Notifier) Notify(ctx context.Context, message string) {
	// TODO: use sync.Cond for performance
	select {
	case <-ctx.Done():
		return
	case notifier.lock <- struct{}{}:
	}

	notifier.writer.WriteString(message)
	notifier.writer.WriteRune('\n')
	notifier.writer.Flush()

	<-notifier.lock
}
