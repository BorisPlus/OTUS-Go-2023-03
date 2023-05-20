package hw05parallelexecution

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatisticsMonitorBasic(t *testing.T) {
	statisticsMonitor := StatisticsMonitor{rwMutex: &sync.RWMutex{}}

	var i uint
	// Статистика в отношении горутин

	statisticsMonitor.SetGoroutinesCountLimit(i)
	require.Equal(t, uint(0), statisticsMonitor.GoroutinesCountLimit(false), "GoroutinesCountLimit set/get. OK.")
	statisticsMonitor.IncGoroutinesCountLimit()
	i++
	require.Equal(t, i, statisticsMonitor.GoroutinesCountLimit(false), "GoroutinesCountLimit increment. OK.")

	i++
	statisticsMonitor.SetGoroutinesCountInit(i)
	require.Equal(t, uint(2), statisticsMonitor.GoroutinesCountInit(), "GoroutinesCountInit set/get. OK.")
	statisticsMonitor.IncGoroutinesCountInit()
	i++
	require.Equal(t, i, statisticsMonitor.GoroutinesCountInit(), "GoroutinesCountInit increment. OK.")

	i++
	statisticsMonitor.SetStartedGoroutinesCount(i)
	require.Equal(t, uint(4), statisticsMonitor.StartedGoroutinesCount(), "StartedGoroutinesCount set/get. OK.")
	statisticsMonitor.IncStartedGoroutinesCount()
	i++
	require.Equal(t, i, statisticsMonitor.StartedGoroutinesCount(), "StartedGoroutinesCount increment. OK.")

	i++
	statisticsMonitor.SetDoneGoroutinesCount(i)
	require.Equal(t, uint(6), statisticsMonitor.DoneGoroutinesCount(), "DoneGoroutinesCount set/get. OK.")
	statisticsMonitor.IncDoneGoroutinesCount()
	i++
	require.Equal(t, i, statisticsMonitor.DoneGoroutinesCount(), "DoneGoroutinesCount increment. OK.")

	// Статистика в отношении задач

	i++
	statisticsMonitor.SetTasksCount(i)
	require.Equal(t, uint(8), statisticsMonitor.TasksCount(), "TasksCount set/get. OK.")
	statisticsMonitor.IncTasksCount()
	i++
	require.Equal(t, i, statisticsMonitor.TasksCount(), "TasksCount increment. OK.")

	i++
	statisticsMonitor.SetTasksCountInit(i)
	require.Equal(t, uint(10), statisticsMonitor.TasksCountInit(), "TasksCountInit set/get. OK.")
	statisticsMonitor.IncTasksCountInit()
	i++
	require.Equal(t, i, statisticsMonitor.TasksCountInit(), "TasksCountInit increment. OK.")

	i++
	statisticsMonitor.SetStartedTasksCount(i)
	require.Equal(t, uint(12), statisticsMonitor.StartedTasksCount(), "StartedTasksCount set/get. OK.")
	statisticsMonitor.IncStartedTasksCount()
	i++
	require.Equal(t, i, statisticsMonitor.StartedTasksCount(), "StartedTasksCount increment. OK.")

	i++
	statisticsMonitor.SetDoneTasksCount(i)
	require.Equal(t, uint(14), statisticsMonitor.DoneTasksCount(), "DoneTasksCount set/get. OK.")
	statisticsMonitor.IncDoneTasksCount()
	i++
	require.Equal(t, i, statisticsMonitor.DoneTasksCount(), "DoneTasksCount increment. OK.")

	i++
	statisticsMonitor.SetErrorsTasksCount(i)
	require.Equal(t, uint(16), statisticsMonitor.ErrorsTasksCount(), "ErrorsTasksCount set/get. OK.")
	statisticsMonitor.IncErrorsTasksCount()
	i++
	require.Equal(t, i, statisticsMonitor.ErrorsTasksCount(), "ErrorsTasksCount increment. OK.")

	i++
	statisticsMonitor.SetErrorsTasksCountLimit(i)
	require.Equal(t, uint(18), statisticsMonitor.ErrorsTasksCountLimit(), "ErrorsTasksCountLimit set/get. OK.")

	// Вычисляемое значение числа горутин, если задач априори меньше

	statisticsMonitor.SetGoroutinesCountLimit(100)
	statisticsMonitor.SetTasksCount(10)
	require.Equal(t, uint(10), statisticsMonitor.GoroutinesCountLimit(true), "GoroutinesCountLimit(true) M<N. OK.")

	statisticsMonitor.SetGoroutinesCountLimit(10)
	statisticsMonitor.SetTasksCount(100)
	require.Equal(t, uint(10), statisticsMonitor.GoroutinesCountLimit(true), "GoroutinesCountLimit(true) M<N. OK.")

	// Критерий останова
	statisticsMonitor.SetErrorsTasksCountLimit(10)
	statisticsMonitor.SetErrorsTasksCount(10)
	require.True(t, statisticsMonitor.DoesErrorsLimitExceeded(), "10 of 10 = ErrorsLimitExceeded. OK.")

	statisticsMonitor.SetErrorsTasksCountLimit(10)
	statisticsMonitor.SetErrorsTasksCount(100)
	require.True(t, statisticsMonitor.DoesErrorsLimitExceeded(), "100 of 10 = ErrorsLimitExceeded. OK.")

	statisticsMonitor.SetErrorsTasksCountLimit(100)
	statisticsMonitor.SetErrorsTasksCount(10)
	require.False(t, statisticsMonitor.DoesErrorsLimitExceeded(), "10 of 100 = so Not ErrorsLimitExceeded. OK.")

	statisticsMonitor.SetErrorsTasksCountLimit(0)
	statisticsMonitor.SetErrorsTasksCount(10)
	require.False(t, statisticsMonitor.DoesErrorsLimitExceeded(), "10 of 0 = so Not ErrorsLimitExceeded. OK.")
}
