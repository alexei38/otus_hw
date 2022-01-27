package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded   = errors.New("errors limit exceeded")
	ErrErrorsNegativeWorkers = errors.New("workers should be > 0")
	ErrErrorsNegativeErrors  = errors.New("errors should be >= 0")
)

type TaskExecutor struct {
	workers   int
	errors    int32
	maxErrors int32
}

type Task func() error

func (t *TaskExecutor) runWorker(wg *sync.WaitGroup, taskCh <-chan Task) {
	defer wg.Done()
	for {
		task, ok := <-taskCh
		if !ok {
			// выходим, если все обработали и канал закрыт
			return
		}
		if err := task(); err != nil {
			atomic.AddInt32(&t.errors, 1)
		}
	}
}

func (t *TaskExecutor) Run(tasks []Task) error {
	taskCh := make(chan Task)
	wg := &sync.WaitGroup{}
	wg.Add(t.workers)
	for i := 0; i < t.workers; i++ {
		go t.runWorker(wg, taskCh)
	}
	for _, task := range tasks {
		if atomic.LoadInt32(&t.errors) >= t.maxErrors {
			break
		}
		taskCh <- task
	}
	// Как записали все задачи в канал, закрываем канал и ждем, когда воркеры их обработают
	close(taskCh)
	wg.Wait()
	if t.errors >= t.maxErrors {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrErrorsNegativeWorkers
	}
	if m < 0 {
		return ErrErrorsNegativeErrors
	}
	executor := &TaskExecutor{
		maxErrors: int32(m),
		workers:   n,
	}
	return executor.Run(tasks)
}
