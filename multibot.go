package multibot

import (
	"context"
	"sync"

	"github.com/cekc/multibot/internal/signal"
)

type Multibot struct {
	fetchers []Fetcher
	handlers []Handler
}

func New() *Multibot {
	var multibot Multibot

	return &multibot
}

func (multibot *Multibot) RegisterFetchers(fetchers ...Fetcher) {
	multibot.fetchers = append(multibot.fetchers, fetchers...)
}

func (multibot *Multibot) RegisterHandlers(handlers ...Handler) {
	multibot.handlers = append(multibot.handlers, handlers...)
}

func (multibot *Multibot) Fetch(ctx context.Context) <-chan Update {
	var channels []<-chan Update
	for _, fetcher := range multibot.fetchers {
		channels = append(channels, fetcher.Fetch(ctx))
	}

	return merge(channels...)
}

func (multibot *Multibot) Serve() {
	multibot.ServeInContext(context.Background())
}

func (multibot *Multibot) ServeInContext(ctx context.Context) {
	ctx = signal.ListenQuit(ctx)

	var wg sync.WaitGroup
	for update := range multibot.Fetch(ctx) {
		for _, handler := range multibot.handlers {
			wg.Add(1)

			go func(handler Handler) {
				defer wg.Done()
				handler.Handle(ctx, update)
			}(handler)
		}
	}

	wg.Wait()
}
