package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestNegativeRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	tasksCount := 50
	tasks := make([]Task, 0, tasksCount)

	var runTasksCount int32

	for i := 0; i < tasksCount; i++ {
		tasks = append(tasks, func() error {
			atomic.AddInt32(&runTasksCount, 1)
			return nil
		})
	}

	tests := []struct {
		workers   int
		maxErrors int
		want      error
	}{
		{workers: -1, maxErrors: 0, want: ErrErrorsNegativeWorkers},
		{workers: 0, maxErrors: 20, want: ErrErrorsNegativeWorkers},
		{workers: 1, maxErrors: -1, want: ErrErrorsNegativeErrors},
	}
	for _, tc := range tests {
		tc := tc
		name := fmt.Sprintf("Test Workers=%d MaxErrors=%d", tc.workers, tc.maxErrors)
		t.Run(name, func(t *testing.T) {
			err := Run(tasks, tc.workers, tc.maxErrors)
			require.Truef(t, errors.Is(err, tc.want), "actual err - %v", err)
		})
	}
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("tasks=0 n=10 m=5", func(t *testing.T) {
		tasksCount := 0
		tasks := make([]Task, 0, tasksCount)

		workersCount := 10
		maxErrorsCount := 5
		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err, "actual err - %v", err)
	})

	t.Run("tasks=50 n=10 m=0", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount), "extra tasks were started")
	})

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
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

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		errorCh := make(chan error)

		mock := clock.NewMock()
		var runPreSleepTasksCount int32
		var runPostSleepTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runPreSleepTasksCount, 1)
				mock.Sleep(time.Second)
				atomic.AddInt32(&runPostSleepTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		go func() {
			err := Run(tasks, workersCount, maxErrorsCount)
			errorCh <- err
			close(errorCh)
		}()

		var runPreSleepTasks int
		for {
			runPreSleepTasks += workersCount
			require.Eventually(t, func() bool {
				return int32(runPreSleepTasks) == atomic.LoadInt32(&runPreSleepTasksCount)
			}, time.Second, 10*time.Millisecond)
			// т.к работают только workersCount, то остальные задачи подвисают, нужно разблокировать Sleep
			mock.Add(time.Second)
			if runPreSleepTasks == tasksCount {
				break
			}
		}
		require.Equal(t, int32(tasksCount), runPreSleepTasksCount)

		require.Eventually(t, func() bool {
			return int32(tasksCount) == atomic.LoadInt32(&runPostSleepTasksCount)
		}, time.Second, 10*time.Millisecond)

		require.Equal(t, int32(tasksCount), runPostSleepTasksCount)

		var err error
		require.Eventually(t, func() bool {
			select {
			case err = <-errorCh:
				return true
			default:
				return false
			}
		}, time.Second, 10*time.Millisecond)
		require.NoError(t, err)
	})
}
