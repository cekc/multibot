package multibot

import (
	"context"
	"sync"
)

func merge(ctx context.Context, channels ...<-chan Update) <-chan Update {
	var wg sync.WaitGroup

	output := make(chan Update)
	done := ctx.Done()

	retransmit := func(input <-chan Update) {
		defer wg.Done()

		for update := range input {
			select {
			case output <- update:

			case <-done:
				return
			}
		}
	}

	for _, input := range channels {
		wg.Add(1)
		go retransmit(input)
	}

	go func() {
		defer close(output)

		wg.Wait()
	}()

	return output
}
