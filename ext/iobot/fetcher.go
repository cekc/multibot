package iobot

import (
	"context"
	"io"
	"os"

	"github.com/cekc/multibot"
)

type Fetcher struct {
	reader   io.Reader
	notifier *Notifier
}

func NewFetcher(reader io.Reader, writer io.Writer) *Fetcher {
	return &Fetcher{reader, NewNotifier(writer)}
}

func NewConsoleFetcher() *Fetcher {
	return NewFetcher(os.Stdin, os.Stdout)
}

func (fetcher *Fetcher) Fetch(ctx context.Context) <-chan multibot.Update {
	updates := make(chan multibot.Update)
	lines := readlines(fetcher.reader)

	go func() {
		defer close(updates)

		for {
			select {
			case <-ctx.Done():
				return

			case line, ok := <-lines:
				if !ok {
					return
				}
				updates <- Update{
					Text:     line,
					Notifier: fetcher.notifier,
				}
			}
		}
	}()

	return updates
}
