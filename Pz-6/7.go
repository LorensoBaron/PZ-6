package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Product struct {
	ID       int
	Name     string
	Price    float64
	Quantity int
}

type Store struct {
	mu       sync.RWMutex
	products map[int]*Product
	sales    int
	revenue  float64
}

func NewStore() *Store {
	return &Store{
		products: make(map[int]*Product),
	}
}

func (s *Store) Add(id int, name string, price float64, qty int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p, ok := s.products[id]; ok {
		p.Quantity += qty
		p.Price = price
		fmt.Printf("Поступление: %s +%d шт. (всего: %d)\n", name, qty, p.Quantity)
	} else {
		s.products[id] = &Product{id, name, price, qty}
		fmt.Printf("Новый товар: %s, %d шт., %.0f руб.\n", name, qty, price)
	}
}

func (s *Store) Sell(id int, qty int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.products[id]
	if !ok || p.Quantity < qty {
		fmt.Printf("Продажа: %s - недостаточно (есть: %d, нужно: %d)\n", p.Name, p.Quantity, qty)
		return false
	}

	p.Quantity -= qty
	revenue := float64(qty) * p.Price
	s.sales += qty
	s.revenue += revenue

	fmt.Printf("Продажа: %s -%d шт. (осталось: %d), сумма: %.0f руб.\n",
		p.Name, qty, p.Quantity, revenue)
	return true
}

func (s *Store) Report() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Println("\nОСТАТКИ: ")
	for _, p := range s.products {
		fmt.Printf("%s: %d шт. x %.0f руб.\n", p.Name, p.Quantity, p.Price)
	}
	fmt.Printf("Продано: %d шт., Выручка: %.0f руб.\n", s.sales, s.revenue)

}

func Supplier(store *Store, id int, wg *sync.WaitGroup, done chan bool) {
	defer wg.Done()
	products := []struct {
		id    int
		name  string
		price float64
	}{
		{1, "Ноутбук", 50000}, {2, "Смартфон", 30000}, {3, "Наушники", 3000},
	}
	for {
		select {
		case <-done:
			return
		default:
			p := products[rand.Intn(len(products))]
			store.Add(p.id, p.name, p.price, rand.Intn(10)+5)
			time.Sleep(time.Millisecond * time.Duration(200+rand.Intn(300)))
		}
	}
}

func Customer(store *Store, id int, wg *sync.WaitGroup, done chan bool) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		default:
			store.Sell(rand.Intn(3)+1, rand.Intn(2)+1)
			time.Sleep(time.Millisecond * time.Duration(100+rand.Intn(200)))
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	store := NewStore()

	store.Add(1, "Ноутбук", 50000, 10)
	store.Add(2, "Смартфон", 30000, 15)
	store.Add(3, "Наушники", 3000, 20)

	var wg sync.WaitGroup
	done := make(chan bool)

	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go Supplier(store, i, &wg, done)
	}

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go Customer(store, i, &wg, done)
	}

	time.Sleep(3 * time.Second)
	close(done)
	wg.Wait()

	store.Report()
}
