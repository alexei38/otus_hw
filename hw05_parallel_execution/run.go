package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Executor interface {
	Run([]Task) error
}

type TaskExecutor struct {
	Workers   int
	errors    int32
	maxErrors int32
}

type Task func() error

func (t *TaskExecutor) runWorker(wg *sync.WaitGroup, taskCh <-chan Task) {
	defer wg.Done()
	for {
		// если набралось достаточное количество ошибок - выходим из воркера
		if atomic.LoadInt32(&t.errors) >= atomic.LoadInt32(&t.maxErrors) {
			return
		}
		task, ok := <-taskCh
		if !ok {
			// выходим, если все обработали и канал закрыт
			return
		}
		err := task()
		if err != nil {
			atomic.AddInt32(&t.errors, 1)
		}
	}
}

func (t *TaskExecutor) getErrors() error {
	if t.errors >= t.maxErrors {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func (t *TaskExecutor) Run(tasks []Task) error {
	// Буферизированный канал, чтобы не блокироваться на записи в канал
	// если все воркеры закончат работу раньше времени из-за ошибок
	taskCh := make(chan Task, len(tasks))
	wg := &sync.WaitGroup{}
	wg.Add(t.Workers)
	for i := 0; i < t.Workers; i++ {
		go t.runWorker(wg, taskCh)
	}
	for _, task := range tasks {
		taskCh <- task
	}
	// Как записали все задачи в канал, закрываем канал и ждем, когда воркеры их обработают
	close(taskCh)
	wg.Wait()
	return t.getErrors()
}

func Run(tasks []Task, n, m int) error {
	var executor Executor = &TaskExecutor{
		maxErrors: int32(m),
		Workers:   n,
	}
	return executor.Run(tasks)
}
