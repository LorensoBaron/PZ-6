package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// SafeLogger потокобезопасный логгер
type SafeLogger struct {
	mu     sync.Mutex
	prefix string
}

// NewSafeLogger создает новый потокобезопасный логгер
func NewSafeLogger(prefix string) *SafeLogger {
	return &SafeLogger{
		prefix: prefix,
	}
}

// Log потокобезопасная запись лога
func (l *SafeLogger) Log(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Форматируем время
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// Создаем префикс с временем
	prefix := fmt.Sprintf("[%s] [%s]", timestamp, l.prefix)

	// Выводим лог
	fmt.Printf(prefix+" "+format+"\n", args...)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	logger := NewSafeLogger("MAIN")

	var wg sync.WaitGroup

	// Запускаем несколько горутин, которые пишут логи
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			workerLogger := NewSafeLogger(fmt.Sprintf("WORKER-%d", id))

			for j := 1; j <= 3; j++ {
				workerLogger.Log("Обработка задачи %d, итерация %d", id, j)
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			}

			workerLogger.Log("Завершил работу")
		}(i)
	}

	logger.Log("Запущено %d воркеров", 5)

	wg.Wait()

	logger.Log("Все воркеры завершили работу")
}
