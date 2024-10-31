package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		wg               sync.WaitGroup
		mu               sync.Mutex
		errCount         int
		errLimitExceeded error
	)

	taskCh := make(chan Task)
	doneCh := make(chan struct{}, 1)

	defer close(doneCh)

	// start producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(taskCh)

		for _, task := range tasks {
			select {
			case <-doneCh:
				return
			default:
				taskCh <- task
			}
		}
	}()

	// start workers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range taskCh {
				err := task()
				if err != nil && m != 0 {
					mu.Lock()
					errCount++

					if errCount == m && errLimitExceeded == nil {
						errLimitExceeded = ErrErrorsLimitExceeded
						doneCh <- struct{}{}
					}
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	return errLimitExceeded
}
