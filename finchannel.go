package main

import (
	"fmt"
	"log"
	"sync"
)

func worker(msg string) (<-chan string, <-chan bool) {
	var wg sync.WaitGroup
	receiver := make(chan string)
	fin := make(chan bool)
	go func() {
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(i int) {
				msg := fmt.Sprintf("%d %s done", i, msg)
				receiver <- msg
				wg.Done()
			}(i)
		}
		wg.Wait()
		fin <- false
	}()
	return receiver, fin
}

func main() {
	receiver, fin := worker("job")
	for {
		select {
		case receive := <-receiver:
			log.Println(receive)
		case <-fin:
			return
		}
	}
}
