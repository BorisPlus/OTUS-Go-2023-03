package hw05parallelexecution

import (
	"errors"
	"sync"
	atomic "sync/atomic"
)

var ErrErrorsLimitExceededWithAtomic = errors.New("errors limit exceeded")

type TaskWithAtomic func() error

func workerWithAtomic(wg *sync.WaitGroup, taskChan <-chan TaskWithAtomic, m int, eCount *int32) {
	defer wg.Done()
	for task := range taskChan {
		if m > 0 && int(atomic.LoadInt32(eCount)) >= m {
			break
		}
		err := task()
		if err != nil {
			atomic.AddInt32(eCount, 1)
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func RunWithAtomic(tasks []TaskWithAtomic, n, m int) error {
	var errorsCount int32

	tasksChan := make(chan TaskWithAtomic)
	go func() {
		defer close(tasksChan)
		for _, t := range tasks {
			if m > 0 && int(atomic.LoadInt32(&errorsCount)) >= m {
				break
			}
			tasksChan <- t
		}
	}()

	wg := sync.WaitGroup{}
	for workerIndex := 1; workerIndex <= n; workerIndex++ {
		wg.Add(1)
		go workerWithAtomic(&wg, tasksChan, m, &errorsCount)
	}

	wg.Wait()

	if m > 0 && int(atomic.LoadInt32(&errorsCount)) >= m {
		return ErrErrorsLimitExceededWithAtomic
	}
	return nil
}
