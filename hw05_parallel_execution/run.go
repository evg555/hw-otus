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
		wg       sync.WaitGroup
		mu       sync.Mutex
		errCount int
	)

	taskCh := make(chan Task)

	// start producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(taskCh)

		for _, task := range tasks {
			mu.Lock()
			if errCount >= m && m != 0 {
				mu.Unlock()
				break
			}
			mu.Unlock()

			taskCh <- task
		}
	}()

	// start workers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range taskCh {
				err := task()
				if err != nil {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	if errCount >= m && m != 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
