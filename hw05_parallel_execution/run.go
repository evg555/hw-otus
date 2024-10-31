package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Counter struct {
	mu       sync.Mutex
	errCount int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	c.errCount++
	c.mu.Unlock()
}

func (c *Counter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.errCount
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		wg         sync.WaitGroup
		errCounter = Counter{}
	)

	taskCh := make(chan Task)

	// start producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(taskCh)

		for _, task := range tasks {
			if errCounter.Get() >= m && m != 0 {
				break
			}

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
					errCounter.Inc()
				}
			}
		}()
	}

	wg.Wait()

	if errCounter.Get() >= m && m != 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
