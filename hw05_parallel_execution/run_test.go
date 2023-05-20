package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRunFirstMTasksErrors(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("If were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})
}

func TestRunAllTasksWithoutAnyError(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("Tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		fmt.Println("One-thread time", int64(sumTime))
		fmt.Println("Multi-thread time", int64(elapsedTime))
		fmt.Println("[One-thread time]/[Multi-thread time]", float32(elapsedTime)/float32(sumTime))

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunWithUnlimitedErrorsCount(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("Unlimited errors count", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			if rand.Intn(2) == 1 {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					return nil
				})
			} else {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					return fmt.Errorf("error from task %d", i)
				})
			}
		}

		workersCount := 5

		maxErrorsCount := 0 // ignore any error count
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Nil(t, err, "actual err - %v", err)

		maxErrorsCount = -1 // ignore 18446744073709551615 errors, brcause -1 is UINT
		err = Run(tasks, workersCount, maxErrorsCount)
		require.Nil(t, err, "actual err - %v", err)
	})
}

func TestRun4TaskWith5Gorutine(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("5 goroutines for 4 tasks", func(t *testing.T) {
		tasksCount := 4
		tasks := make([]Task, 0, tasksCount)
		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 0 // ignore any error count
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Nil(t, err, "actual err - %v", err)
	})
}
