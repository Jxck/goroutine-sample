package main

import (
	"fmt"
	"log"
	"sync"
)

func worker(msg string) <-chan string {
	var wg sync.WaitGroup
	receiver := make(chan string)
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
		close(receiver)
	}()
	return receiver
}

func main() {
	receiver := worker("job")
	for {
		receive, ok := <-receiver
		if !ok {
			log.Println("closed")
			return
		}
		log.Println(receive)
	}
}
