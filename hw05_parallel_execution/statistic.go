package hw05parallelexecution

import (
	"fmt"
	"sync"
)

// StatisticsMonitor - структура учета статистики выполнения набора задач.
type StatisticsMonitor struct {
	rwMutex *sync.RWMutex
	// Статистика в отношении горутин

	goroutinesCountLimit   uint // Лимит на общее число горутин
	goroutinesCountInit    uint // Всего подготавливалось к запуcку горутин
	startedGoroutinesCount uint // Всего было запущено горутин
	doneGoroutinesCount    uint // Всего горутин исполнилось

	// Статистика в отношении задач

	tasksCount            uint // Общее число задач
	tasksCountInit        uint // Всего подготавливалось к запуcку задач
	startedTasksCount     uint // Всего было запущено задач
	doneTasksCount        uint // Всего задач исполнилось успешно
	errorsTasksCountLimit uint // Лимит на число задач, завершившихся с ошибками
	errorsTasksCount      uint // Всего задач завершилось с ошибками
}

// StatisticsMonitor - конструктор с обязательными полями.
func NewStatisticsMonitor(mtx *sync.RWMutex, m, n, tasksCount uint) StatisticsMonitor {
	stat := StatisticsMonitor{}
	stat.rwMutex = mtx
	stat.SetErrorsTasksCountLimit(m)
	stat.SetTasksCount(tasksCount)
	stat.SetGoroutinesCountLimit(n)
	return stat
}

func (statMonitor StatisticsMonitor) String() string {
	return fmt.Sprintf(`СТАТИСТИКА РАБОТЫ
    ГОРУТИНЫ
        Лимит на общее число горутин: %d
        Лимит на общее число горутин (вычилено): %d
        Всего подготавливалось к запуcку горутин: %d
        Всего было запущено горутин: %d
        Всего горутин исполнилось: %d
    ЗАДАЧИ
        Общее число задач: %d
        Всего подготавливалось к запуcку задач: %d
        Всего было запущено задач: %d
        Всего задач исполнилось успешно: %d
        Лимит на число задач, завершившихся с ошибками: %d
        Всего задач завершилось с ошибками: %d`,
		// Статистика в отношении горутин
		statMonitor.GoroutinesCountLimit(false),
		statMonitor.GoroutinesCountLimit(true),
		statMonitor.GoroutinesCountInit(),
		statMonitor.StartedGoroutinesCount(),
		statMonitor.DoneGoroutinesCount(),
		// Статистика в отношении задач
		statMonitor.TasksCount(),
		statMonitor.TasksCountInit(),
		statMonitor.StartedTasksCount(),
		statMonitor.DoneTasksCount(),
		statMonitor.ErrorsTasksCountLimit(),
		statMonitor.ErrorsTasksCount(),
	)
}

// SetErrorsTasksCountLimit() - устаноавливает ограничение на число задач, завершающихся ошибками.
// При этом значение 0 не ограничивает число задач с ошибками.
func (statMonitor *StatisticsMonitor) SetErrorsTasksCountLimit(limit uint) {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.errorsTasksCountLimit = limit
}

// ErrorsTasksCountLimit() - ограничение на число задач, завершающихся ошибками.
// При этом значение 0 не ограничивает число задач с ошибками.
func (statMonitor *StatisticsMonitor) ErrorsTasksCountLimit() uint {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	return statMonitor.errorsTasksCountLimit
}

// cleverGoroutinesCountLimit() - ограничение на число горутин "по-умному".
// И без этого отработает, но если задач меньше, чем воркеров,
// то зачем тогда все воркеры запускать.
// Воркеров столько - сколько и задач, если задач априори меньше воркеров.
func (statMonitor *StatisticsMonitor) cleverGoroutinesCountLimit() uint {
	// TODO: double lock if call from other method
	// defer statMonitor.rwMutex.Unlock()
	// statMonitor.rwMutex.Lock()
	if statMonitor.goroutinesCountLimit <= statMonitor.tasksCount {
		return statMonitor.goroutinesCountLimit
	}
	return statMonitor.tasksCount
}

// DoesErrorsLimitExceeded() - критерий останова.
func (statMonitor *StatisticsMonitor) DoesErrorsLimitExceeded() bool {
	defer statMonitor.rwMutex.RUnlock()
	statMonitor.rwMutex.RLock()
	return statMonitor.errorsTasksCountLimit > 0 && statMonitor.errorsTasksCount >= statMonitor.errorsTasksCountLimit
}

// goroutinesCountLimit

// SetGoroutinesCountLimit() - устанавливает ограничение на число горутин.
func (statMonitor *StatisticsMonitor) SetGoroutinesCountLimit(v uint) {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.goroutinesCountLimit = v
}

func (statMonitor *StatisticsMonitor) GoroutinesCountLimit(clever bool) uint {
	defer statMonitor.rwMutex.RUnlock()
	statMonitor.rwMutex.RLock()
	if clever {
		return statMonitor.cleverGoroutinesCountLimit()
	}
	return statMonitor.goroutinesCountLimit
}

func (statMonitor *StatisticsMonitor) IncGoroutinesCountLimit() {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.goroutinesCountLimit++
}

// goroutinesCountInit

func (statMonitor *StatisticsMonitor) SetGoroutinesCountInit(v uint) {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.goroutinesCountInit = v
}

func (statMonitor *StatisticsMonitor) GoroutinesCountInit() uint {
	defer statMonitor.rwMutex.RUnlock()
	statMonitor.rwMutex.RLock()
	return statMonitor.goroutinesCountInit
}

func (statMonitor *StatisticsMonitor) IncGoroutinesCountInit() {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.goroutinesCountInit++
}

// startedGoroutinesCount

func (statMonitor *StatisticsMonitor) SetStartedGoroutinesCount(v uint) {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.startedGoroutinesCount = v
}

func (statMonitor *StatisticsMonitor) StartedGoroutinesCount() uint {
	defer statMonitor.rwMutex.RUnlock()
	statMonitor.rwMutex.RLock()
	return statMonitor.startedGoroutinesCount
}

func (statMonitor *StatisticsMonitor) IncStartedGoroutinesCount() {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.startedGoroutinesCount++
}

// doneGoroutinesCount

func (statMonitor *StatisticsMonitor) SetDoneGoroutinesCount(v uint) {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.doneGoroutinesCount = v
}

func (statMonitor *StatisticsMonitor) DoneGoroutinesCount() uint {
	defer statMonitor.rwMutex.RUnlock()
	statMonitor.rwMutex.RLock()
	return statMonitor.doneGoroutinesCount
}

func (statMonitor *StatisticsMonitor) IncDoneGoroutinesCount() {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.doneGoroutinesCount++
}

// tasksCount

func (statMonitor *StatisticsMonitor) SetTasksCount(v uint) {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.tasksCount = v
}

func (statMonitor *StatisticsMonitor) TasksCount() uint {
	defer statMonitor.rwMutex.RUnlock()
	statMonitor.rwMutex.RLock()
	return statMonitor.tasksCount
}

func (statMonitor *StatisticsMonitor) IncTasksCount() {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.tasksCount++
}

// tasksCountInit

func (statMonitor *StatisticsMonitor) SetTasksCountInit(v uint) {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.tasksCountInit = v
}

func (statMonitor *StatisticsMonitor) TasksCountInit() uint {
	defer statMonitor.rwMutex.RUnlock()
	statMonitor.rwMutex.RLock()
	return statMonitor.tasksCountInit
}

func (statMonitor *StatisticsMonitor) IncTasksCountInit() {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.tasksCountInit++
}

// startedTasksCount

func (statMonitor *StatisticsMonitor) SetStartedTasksCount(v uint) {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.startedTasksCount = v
}

func (statMonitor *StatisticsMonitor) StartedTasksCount() uint {
	defer statMonitor.rwMutex.RUnlock()
	statMonitor.rwMutex.RLock()
	return statMonitor.startedTasksCount
}

func (statMonitor *StatisticsMonitor) IncStartedTasksCount() {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.startedTasksCount++
}

// doneTasksCount

func (statMonitor *StatisticsMonitor) SetDoneTasksCount(v uint) {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.doneTasksCount = v
}

func (statMonitor *StatisticsMonitor) DoneTasksCount() uint {
	defer statMonitor.rwMutex.RUnlock()
	statMonitor.rwMutex.RLock()
	return statMonitor.doneTasksCount
}

func (statMonitor *StatisticsMonitor) IncDoneTasksCount() {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.doneTasksCount++
}

// SetErrorsTasksCount() - число задач, завершающихся ошибками.
func (statMonitor *StatisticsMonitor) SetErrorsTasksCount(v uint) {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.errorsTasksCount = v
}

func (statMonitor *StatisticsMonitor) ErrorsTasksCount() uint {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	return statMonitor.errorsTasksCount
}

func (statMonitor *StatisticsMonitor) IncErrorsTasksCount() {
	defer statMonitor.rwMutex.Unlock()
	statMonitor.rwMutex.Lock()
	statMonitor.errorsTasksCount++
}
