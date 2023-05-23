package hw05parallelexecution

import (
	"sync"
	atomic "sync/atomic"
)

func workerWithAtomic(wg *sync.WaitGroup, tasksChan <-chan Task, errorsCount *int64, m int) {
	defer wg.Done()
	for task := range tasksChan {
		if m > 0 && *errorsCount >= int64(m) {
			break
		}
		err := task()
		if err != nil {
			atomic.AddInt64(errorsCount, 1)
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func RunWithAtomic(tasks []Task, n, m int) error {
	var errorsCount int64
	tasksChan := make(chan Task)
	go func() {
		defer close(tasksChan)
		for _, t := range tasks {
			if m > 0 && errorsCount >= int64(m) {
				break
			}
			tasksChan <- t
		}
	}()
	wg := sync.WaitGroup{}
	for workerIndex := 1; workerIndex <= n; workerIndex++ {
		wg.Add(1)
		go workerWithAtomic(&wg, tasksChan, &errorsCount, m)
	}
	wg.Wait()

	if m > 0 && errorsCount >= int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
