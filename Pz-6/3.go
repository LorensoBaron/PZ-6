package main

import (
	"fmt"
	"sync"
)

type TaskQueue struct {
	mu    sync.Mutex
	tasks []string
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		tasks: make([]string, 0),
	}
}

func (q *TaskQueue) Add(task string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks = append(q.tasks, task)
	fmt.Printf("Добавлена задача: %s (в очереди: %d)\n", task, len(q.tasks))
}

func (q *TaskQueue) Get() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tasks) == 0 {
		return "", false
	}

	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	fmt.Printf("Выполнена задача: %s (осталось: %d)\n", task, len(q.tasks))
	return task, true
}

func main() {
	queue := NewTaskQueue()
	var wg sync.WaitGroup

	// Добавляем задачи (3 продюсера)
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			queue.Add(fmt.Sprintf("Задача #%d", id))
		}(i)
	}

	// Обрабатываем задачи (2 воркера)
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				task, ok := queue.Get()
				if !ok {
					break
				}
				_ = task
			}
		}(i)
	}

	wg.Wait()
}
