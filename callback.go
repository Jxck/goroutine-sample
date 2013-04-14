package main

import (
	"log"
	"time"
)

func async(cb func(ch chan string)) {
	ch := make(chan string)
	go cb(ch)
	ch <- "hello"
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go async(func(ch chan string) {
		log.Println(<-ch)
		wg.Done()
	})
	wg.Wait()
}
