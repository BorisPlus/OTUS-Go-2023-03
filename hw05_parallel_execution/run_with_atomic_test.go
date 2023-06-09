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

func TestRunWithAtomicFirstMTasksErrors(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("If were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]TaskWithAtomic, 0, tasksCount)

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
		err := RunWithAtomic(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceededWithAtomic), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})
}

func TestRunWithAtomicAllTasksWithoutAnyError(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("Tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]TaskWithAtomic, 0, tasksCount)

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
		err := RunWithAtomic(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunWithAtomicWithUnlimitedErrorsCount(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("Unlimited errors count", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]TaskWithAtomic, 0, tasksCount)

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
		err := RunWithAtomic(tasks, workersCount, maxErrorsCount)
		require.Nil(t, err, "actual err - %v", err)

		maxErrorsCount = -1 // ignore 18446744073709551615 errors, brcause -1 is UINT
		err = RunWithAtomic(tasks, workersCount, maxErrorsCount)
		require.Nil(t, err, "actual err - %v", err)
	})
}
