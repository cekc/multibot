package multibot

import (
	"context"
	"errors"

	"github.com/cekc/multibot/worker"
)

type Multibot struct {
	WorkerPool WorkerPool

	fetchers            []Fetcher
	handlers            []Handler
	cancelProcessingCtx context.CancelFunc
	ranOutOfUpdates     chan struct{}
}

func New() *Multibot {
	return &Multibot{
		WorkerPool: worker.NewPool(),
	}
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

func (multibot *Multibot) Process() {
	multibot.ranOutOfUpdates = make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	multibot.cancelProcessingCtx = cancel

	for update := range multibot.fetch(ctx) {
		for _, handler := range multibot.handlers {
			handler, update := handler, update
			multibot.WorkerPool.Submit(func() {
				handler.Handle(ctx, update)
			})
		}
	}

	close(multibot.ranOutOfUpdates)
}

func (multibot *Multibot) RanOutOfUpdates() <-chan struct{} {
	return multibot.ranOutOfUpdates
}

func (multibot *Multibot) Shutdown(ctx context.Context) error {
	multibot.cancelProcessingCtx()

	workersDone := make(chan struct{})
	go func() {
		defer close(workersDone)
		multibot.WorkerPool.Wait()
	}()

	select {
	case <-ctx.Done():
		return errors.New("Multibot shutdown: context expired")
	case <-workersDone:
	}

	return nil
}
