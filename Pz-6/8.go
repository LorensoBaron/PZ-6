package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job представляет работу для выполнения
type Job struct {
	ID   int
	Data string
}

// JobResult представляет результат работы
type JobResult struct {
	JobID    int
	WorkerID int
	Success  bool
	Message  string
}

// WorkerPool управляет пулом воркеров
type WorkerPool struct {
	jobQueue    chan Job
	resultQueue chan JobResult
	workers     int
	stats       struct {
		sync.Mutex
		processed int
		failed    int
		byWorker  map[int]int
	}
	wg sync.WaitGroup
}

// NewWorkerPool создает новый пул воркеров
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	return &WorkerPool{
		jobQueue:    make(chan Job, queueSize),
		resultQueue: make(chan JobResult, queueSize),
		workers:     workers,
		stats: struct {
			sync.Mutex
			processed int
			failed    int
			byWorker  map[int]int
		}{
			byWorker: make(map[int]int),
		},
	}
}

// Start запускает воркеров
func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker выполняет обработку задач
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for job := range wp.jobQueue {
		// Имитация обработки
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

		// Случайные ошибки для демонстрации
		success := rand.Float32() > 0.1

		result := JobResult{
			JobID:    job.ID,
			WorkerID: id,
			Success:  success,
			Message:  fmt.Sprintf("Обработка: %s", job.Data),
		}

		wp.resultQueue <- result

		// Обновление статистики
		wp.stats.Lock()
		wp.stats.processed++
		wp.stats.byWorker[id]++
		if !success {
			wp.stats.failed++
		}
		wp.stats.Unlock()
	}
}

// AddJob добавляет задачу в очередь
func (wp *WorkerPool) AddJob(job Job) {
	wp.jobQueue <- job
}

// GetResults собирает результаты
func (wp *WorkerPool) GetResults() []JobResult {
	var results []JobResult

	// Закрываем jobQueue после завершения добавления задач
	// и ожидаем завершения воркеров
	go func() {
		wp.wg.Wait()
		close(wp.resultQueue)
	}()

	for result := range wp.resultQueue {
		results = append(results, result)
	}

	return results
}

// Stop останавливает пул
func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
}

// PrintStats выводит статистику
func (wp *WorkerPool) PrintStats() {
	wp.stats.Lock()
	defer wp.stats.Unlock()

	fmt.Printf("\n=== Статистика ===\n")
	fmt.Printf("Обработано задач: %d\n", wp.stats.processed)
	fmt.Printf("Ошибок: %d\n", wp.stats.failed)
	fmt.Printf("Успешно: %d\n", wp.stats.processed-wp.stats.failed)

	fmt.Println("\nНагрузка на воркеров:")
	for id, count := range wp.stats.byWorker {
		fmt.Printf("  Воркер %d: %d задач\n", id, count)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Создание пула воркеров
	pool := NewWorkerPool(5, 20)

	// Запуск воркеров
	pool.Start()

	// Продюсеры (отправка задач)
	var producersWG sync.WaitGroup

	for p := 1; p <= 3; p++ {
		producersWG.Add(1)
		go func(producerID int) {
			defer producersWG.Done()

			for i := 1; i <= 10; i++ {
				job := Job{
					ID:   producerID*1000 + i,
					Data: fmt.Sprintf("Задача %d от продюсера %d", i, producerID),
				}
				pool.AddJob(job)
				fmt.Printf("Продюсер %d: добавлена задача %d\n", producerID, job.ID)
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(20)))
			}
		}(p)
	}

	// Ожидание завершения продюсеров
	producersWG.Wait()

	// Остановка пула (закрытие очереди задач)
	pool.Stop()

	// Получение результатов
	results := pool.GetResults()

	// Вывод результатов
	fmt.Printf("\n=== Результаты (%d) ===\n", len(results))

	// ИСПРАВЛЕНИЕ: используем переменную i в теле цикла
	limit := 10
	if len(results) < limit {
		limit = len(results)
	}

	for i := 0; i < limit; i++ {
		result := results[i]
		status := "✓"
		if !result.Success {
			status = "✗"
		}
		fmt.Printf("%s Задача %d (воркер %d): %s\n",
			status, result.JobID, result.WorkerID, result.Message)
	}

	// Вывод статистики
	pool.PrintStats()
}
