package multibot

import (
	"context"
	"sync"
)

type Multibot struct {
	fetchers []Fetcher
	handlers []Handler
}

func New() *Multibot {
	var multibot Multibot

	return &multibot
}

func (multibot *Multibot) AddFetcher(fetcher Fetcher) {
	multibot.fetchers = append(multibot.fetchers, fetcher)
}

func (multibot *Multibot) AddHandler(handler Handler) {
	multibot.handlers = append(multibot.handlers, handler)
}

func (multibot *Multibot) fetch(ctx context.Context) <-chan Update {
	var channels []<-chan Update
	for _, fetcher := range multibot.fetchers {
		channels = append(channels, fetcher.Fetch(ctx))
	}

	return merge(channels...)
}

func (multibot *Multibot) Process(ctx context.Context) {
	var wg sync.WaitGroup
	for update := range multibot.fetch(ctx) {
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
