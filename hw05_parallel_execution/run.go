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
	if m == 0 {
		return ErrErrorsLimitExceeded
	}

	var (
		wg                    sync.WaitGroup
		errorLimit            = int32(m)
		isErrorsLimitExceeded bool
	)

	taskChan := make(chan Task)
	quitChan := make(chan struct{})

	// Start exactly n workers.
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			for task := range taskChan {
				if err := task(); err != nil {
					if atomic.AddInt32(&errorLimit, -1) == 0 {
						quitChan <- struct{}{}
						isErrorsLimitExceeded = true
					}
				}
			}
			wg.Done()
		}()
	}

	// Feed workers while errors limit is not exceeded.
	go feed(tasks, taskChan, quitChan)

	wg.Wait()
	close(quitChan)

	if isErrorsLimitExceeded {
		return ErrErrorsLimitExceeded
	}

	return nil
}

// Returns true if all tasks were successfully fed to taskChannel.
func feed(tasks []Task, taskChannel chan<- Task, quitChan <-chan struct{}) {
	defer func() {
		close(taskChannel)
		for range quitChan {
		}
	}()

	for _, task := range tasks {
		select {
		case <-quitChan:
			return
		default:
		}

		select {
		case <-quitChan:
			return
		case taskChannel <- task:
			// Next task was dispatched.
		}
	}
}
