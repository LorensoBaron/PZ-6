package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Cinema struct {
	mu    sync.Mutex
	seats [38]bool
}

func NewCinema() *Cinema {
	return &Cinema{}
}

func (c *Cinema) Book(seatNum int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if seatNum < 1 || seatNum > 38 {
		fmt.Printf("Место %d не существует\n", seatNum)
		return false
	}

	if c.seats[seatNum-1] {
		fmt.Printf("Место %d уже занято\n", seatNum)
		return false
	}

	c.seats[seatNum-1] = true
	fmt.Printf("Место %d успешно забронировано\n", seatNum)
	return true
}

func (c *Cinema) ShowSeats() {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println("\nСхема зала (X - занято):")
	for i := 0; i < 38; i++ {
		if i%10 == 0 && i != 0 {
			fmt.Println()
		}
		if c.seats[i] {
			fmt.Printf("[ X]")
		} else {
			fmt.Printf("[%2d]", i+1)
		}
	}
	fmt.Println("\n")
}

func main() {
	cinema := NewCinema()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Система бронирования мест в кинотеатре")
	fmt.Println("Команды: book N - забронировать место N, show - показать зал, exit - выход")

	for {
		fmt.Print("> ")
		scanner.Scan()
		cmd := scanner.Text()
		parts := strings.Fields(cmd)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "book":
			if len(parts) != 2 {
				fmt.Println("Использование: book N")
				continue
			}
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Введите число")
				continue
			}
			// Запускаем бронь в горутине для имитации одновременных запросов
			go cinema.Book(num)

		case "show":
			cinema.ShowSeats()

		case "exit":
			fmt.Println("До свидания!")
			return

		default:
			fmt.Println("Неизвестная команда")
		}
	}
}
