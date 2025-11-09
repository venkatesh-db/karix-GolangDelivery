package main

import (
	"fmt"
	"sync"
)

type Reelstats struct {
	views    int
	likes    int
	comments int
	shared   int
	mu       sync.Mutex
}

// one post in instgram
// background 3 go routines view likes comments

func simulateviews(p *Reelstats, wg *sync.WaitGroup) {
 defer wg.Done()

	p.mu.Lock()
	fmt.Println("simulateviews")
	p.views=1000
	p.mu.Unlock()

}

func simulatelikes(p *Reelstats, wg *sync.WaitGroup) {
 defer wg.Done()
 	p.mu.Lock()
	fmt.Println("simulatelikes")
	p.likes=550
		p.mu.Unlock()
}

func simulatecomments(p *Reelstats, wg *sync.WaitGroup) {
 defer wg.Done()
	p.mu.Lock()
	fmt.Println("simulatecomments")
	p.comments=5
	p.mu.Unlock()
}

func main() {

	var ajith Reelstats
	var wg sync.WaitGroup

	wg.Add(3)

	go simulateviews(&ajith, &wg)
	go simulatelikes(&ajith, &wg)
	go simulatecomments(&ajith, &wg)

	wg.Wait() // wait till 3 gororutien compleets

	fmt.Println(ajith)
}
