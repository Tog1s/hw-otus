package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var waitGroup sync.WaitGroup
	var errorLimit int32

	channel := make(chan Task, len(tasks))
	for _, task := range tasks {
		channel <- task

	}
	close(channel)

	waitGroup.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer waitGroup.Done()
			for task := range channel {
				if atomic.LoadInt32(&errorLimit) >= int32(m) {
					return
				}
				if task() != nil {
					atomic.AddInt32(&errorLimit, 1)
				}
			}
		}()
	}

	waitGroup.Wait()
	if errorLimit >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
