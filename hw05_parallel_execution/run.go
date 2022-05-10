package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m == 0 {
		return ErrErrorsLimitExceeded
	}

	var wg sync.WaitGroup

	taskChan := make(chan Task)
	taskResultChan := make(chan error)
	quitChan := make(chan struct{})

	// Start exactly n workers.
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			for task := range taskChan {
				taskResultChan <- task()
			}
			wg.Done()
		}()
	}

	// Start async error counter.
	go errorLimiter(taskResultChan, m, quitChan)

	// Feed workers while errors limit is not exceeded.
	isAllTasksDispatched := feed(tasks, taskChan, quitChan)

	// Finish work.
	wg.Wait()
	close(taskResultChan)

	if !isAllTasksDispatched {
		return ErrErrorsLimitExceeded
	}

	return nil
}

// Writes to quitChan if errLimit of errors from taskResultChan exceeded.
func errorLimiter(taskResultChan <-chan error, errLimit int, quitChan chan<- struct{}) {
	for err := range taskResultChan {
		if err == nil {
			continue
		}

		errLimit--
		if errLimit == 0 {
			quitChan <- struct{}{}
		}
	}
}

// Returns true if all tasks were successfully fed to taskChannel.
func feed(tasks []Task, taskChannel chan<- Task, quitChan <-chan struct{}) bool {
	defer close(taskChannel)

	for _, task := range tasks {
		select {
		case <-quitChan:
			return false
		default:
		}

		select {
		case <-quitChan:
			return false
		case taskChannel <- task:
			// Next task was dispatched.
		}
	}

	return true
}
