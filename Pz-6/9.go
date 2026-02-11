package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Metrics struct {
	mu           sync.RWMutex
	success      int
	failed       int
	totalTime    int
	requestCount int
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RecordSuccess(duration int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.success++
	m.totalTime += duration
	m.requestCount++
}

func (m *Metrics) RecordFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failed++
}

func (m *Metrics) Report() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fmt.Println("\n=== ТЕКУЩИЕ МЕТРИКИ ===")
	fmt.Printf("Успешных запросов: %d\n", m.success)
	fmt.Printf("Неуспешных запросов: %d\n", m.failed)
	fmt.Printf("Всего запросов: %d\n", m.success+m.failed)

	if m.requestCount > 0 {
		avgTime := m.totalTime / m.requestCount
		fmt.Printf("Среднее время ответа: %d мс\n", avgTime)
	}
	fmt.Println("========================\n")
}

func main() {
	metrics := NewMetrics()
	scanner := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup

	fmt.Println("Сбор метрик запросов")
	fmt.Println("Команды: success N - успешный запрос (N мс), failed - неуспешный, report - показать метрики, exit - выход")

	for {
		fmt.Print("> ")
		scanner.Scan()
		cmd := scanner.Text()
		parts := strings.Fields(cmd)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "success":
			duration := 100
			if len(parts) == 2 {
				d, err := strconv.Atoi(parts[1])
				if err == nil {
					duration = d
				}
			}

			wg.Add(1)
			go func(d int) {
				defer wg.Done()
				time.Sleep(time.Millisecond * time.Duration(d))
				metrics.RecordSuccess(d)
				fmt.Printf("✓ Успешный запрос (%d мс)\n", d)
			}(duration)

		case "failed":
			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(time.Millisecond * 50)
				metrics.RecordFailed()
				fmt.Printf("✗ Неуспешный запрос\n")
			}()

		case "report":
			metrics.Report()

		case "exit":
			wg.Wait()
			fmt.Println("Финальный отчет:")
			metrics.Report()
			fmt.Println("До свидания!")
			return

		default:
			fmt.Println("Неизвестная команда")
		}
	}
}
