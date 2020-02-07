package filebot

import (
	"context"
	"io"
	"os"

	"github.com/cekc/multibot"
)

type Fetcher struct {
	reader io.Reader
}

func NewFetcher(reader io.Reader) *Fetcher {
	return &Fetcher{reader}
}

func NewConsoleFetcher() *Fetcher {
	return NewFetcher(os.Stdin)
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
				updates <- multibot.Update{
					Body: line,
					From: "mysterious old man",
				}
			}
		}
	}()

	return updates
}
