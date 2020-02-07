package worker

import (
	"sync"
)

type Pool struct {
	wg sync.WaitGroup
}

func NewPool() *Pool {
	return &Pool{}
}

func (pool *Pool) Submit(task func()) {
	pool.wg.Add(1)
	go func() {
		defer pool.wg.Done()
		task()
	}()
}

func (pool *Pool) Wait() {
	pool.wg.Wait()
}
