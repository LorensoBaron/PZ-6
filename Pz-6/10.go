package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Pipeline struct {
	mu     sync.Mutex
	errors []string
}

func (p *Pipeline) addError(err string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.errors = append(p.errors, err)
}

func (p *Pipeline) printErrors() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.errors) == 0 {
		fmt.Println("Ошибок нет")
		return
	}

	fmt.Println("\nОшибки в конвейере:")
	for i, e := range p.errors {
		fmt.Printf("%d. %s\n", i+1, e)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	p := &Pipeline{
		errors: make([]string, 0),
	}

	fmt.Println("Конвейер обработки чисел")
	fmt.Println("Введите числа через пробел:")
	scanner.Scan()
	input := scanner.Text()
	nums := strings.Fields(input)

	var wg sync.WaitGroup

	for _, numStr := range nums {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Printf("Пропускаем '%s' - не число\n", numStr)
			continue
		}

		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			// Каналы для одной цепочки обработки
			ch1 := make(chan int, 1)
			ch2 := make(chan int, 1)
			ch3 := make(chan int, 1)

			// Стадия 1
			go func() {
				if n < 0 {
					p.addError(fmt.Sprintf("Стадия 1: число %d отрицательное", n))
					close(ch1)
					return
				}
				ch1 <- n * 2
				close(ch1)
			}()

			// Стадия 2
			go func() {
				val, ok := <-ch1
				if !ok {
					close(ch2)
					return
				}
				if val > 100 {
					p.addError(fmt.Sprintf("Стадия 2: число %d больше 100", val))
				}
				ch2 <- val
				close(ch2)
			}()

			// Стадия 3
			go func() {
				val, ok := <-ch2
				if !ok {
					close(ch3)
					return
				}
				ch3 <- val + 5
				close(ch3)
			}()

			// Стадия 4
			go func() {
				val, ok := <-ch3
				if !ok {
					return
				}
				if val%10 == 0 {
					p.addError(fmt.Sprintf("Стадия 4: число %d делится на 10", val))
				}
				fmt.Printf("Результат: %d\n", val)
			}()
		}(num)
	}

	wg.Wait()
	p.printErrors()
}
