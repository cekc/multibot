package multibot

import (
	"sync"
)

func merge(channels ...<-chan Update) <-chan Update {
	var wg sync.WaitGroup

	output := make(chan Update)
	retransmit := func(input <-chan Update) {
		defer wg.Done()

		for update := range input {
			output <- update
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
