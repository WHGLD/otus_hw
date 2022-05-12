package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run n - кол-во горутин
// Run m - кол-во ошибок
func Run(tasks []Task, n, m int) error {
	if n <= 0 || m <= 0 {
		return ErrErrorsLimitExceeded // todo unit test
	}

	// спорный момент, возможно не надо так делать и в тестах упадет
	if len(tasks) < n {
		n = len(tasks)
	}

	wg := sync.WaitGroup{}
	tasksChan := make(chan Task)

	var errsCount int32

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// разбираем задачи из потока и обновляем ErrsCount
			for task := range tasksChan {
				result := task()
				if result != nil {
					atomic.AddInt32(&errsCount, 1)
				}
			}
		}()
	}

	var response error

	// пишем задачи в поток
	for _, task := range tasks {
		tasksChan <- task

		if int(atomic.LoadInt32(&errsCount)) >= m {
			response = ErrErrorsLimitExceeded
			break
		}
	}

	close(tasksChan)
	wg.Wait()

	return response
}
