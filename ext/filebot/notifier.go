package filebot

import (
	"bufio"
	"context"
	"io"
	"sync"
)

type Notifier struct {
	writer bufio.Writer
	mut    sync.Mutex
}

func CreateNotifier(writer io.Writer) Notifier {
	return Notifier{
		writer: *bufio.NewWriter(writer),
	}
}

func (notifier *Notifier) Notify(ctx context.Context, message string) {
	defer notifier.mut.Unlock()
	notifier.mut.Lock()

	notifier.writer.WriteString(message)
	notifier.writer.WriteRune('\n')
	notifier.writer.Flush()
}
