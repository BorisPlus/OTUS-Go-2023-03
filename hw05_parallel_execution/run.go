package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(wg *sync.WaitGroup, tasksChan <-chan Task, stat *StatisticsMonitor) {
	defer func(stat *StatisticsMonitor) {
		stat.IncDoneGoroutinesCount()
		wg.Done()
	}(stat)

	stat.IncStartedGoroutinesCount()

	for task := range tasksChan {
		stat.IncTasksCountInit()
		if stat.DoesErrorsLimitExceeded() {
			break
		}
		stat.IncStartedTasksCount()
		err := task()
		if err != nil {
			stat.IncErrorsTasksCount()
		} else {
			stat.IncDoneTasksCount()
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workTogetherTasksCountLimit, errorsCountLimit int) error {
	mtx := sync.RWMutex{}

	stat := NewStatisticsMonitor(&mtx, uint(errorsCountLimit), uint(workTogetherTasksCountLimit), uint(len(tasks)))
	fmt.Printf("\nИСХОДНАЯ\n%s\n", stat)

	defer func() {
		fmt.Printf("\nИТОГОВАЯ\n%s\n", stat)
	}()

	// COMMIT_REQUEST 
	tasksChan := make(chan Task, len(tasks))
	// tasksChan := make(chan Task)
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)

	var workerIndex uint
	wg := sync.WaitGroup{}
	for workerIndex = 1; workerIndex <= stat.GoroutinesCountLimit(true); workerIndex++ {
		stat.IncGoroutinesCountInit()

		wg.Add(1)
		go worker(&wg, tasksChan, &stat)
	}

	wg.Wait()

	if stat.DoesErrorsLimitExceeded() {
		fmt.Println("Errors was limit!!!")
		return ErrErrorsLimitExceeded
	}
	return nil
}
