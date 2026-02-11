package main

import (
	"fmt"
	"sync"
)

type VoteCounter struct {
	mu    sync.Mutex
	votes map[string]int
}

func NewVoteCounter() *VoteCounter {
	return &VoteCounter{
		votes: make(map[string]int),
	}
}

func (v *VoteCounter) Vote(candidate string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.votes[candidate]++
	fmt.Printf("Голос за %s (всего: %d)\n", candidate, v.votes[candidate])
}

func (v *VoteCounter) Results() {
	v.mu.Lock()
	defer v.mu.Unlock()

	fmt.Println("\n=== ИТОГИ ГОЛОСОВАНИЯ ===")
	fmt.Println("==========================")
	total := 0
	for candidate, votes := range v.votes {
		fmt.Printf("%s: %d голосов\n", candidate, votes)
		total += votes
	}
	fmt.Println("==========================")
	fmt.Printf("ВСЕГО: %d голосов\n", total)
}

func main() {
	counter := NewVoteCounter()
	var wg sync.WaitGroup

	candidates := []string{"Кандидат А", "Кандидат Б", "Кандидат В"}

	// 100 избирателей голосуют
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			candidate := candidates[n%len(candidates)]
			counter.Vote(candidate)
		}(i)
	}

	wg.Wait()
	counter.Results()
}
