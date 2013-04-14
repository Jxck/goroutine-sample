package main

import (
	"fmt"
	"log"
)

func worker(msg string) <-chan string {
	receiver := make(chan string)
	for i := 0; i < 3; i++ {
		go func(i int) {
			msg := fmt.Sprintf("%d %s done", i, msg)
			receiver <- msg
		}(i)
	}
	return receiver
}

func main() {
	receiver := worker("job")
	for i := 0; i < 3; i++ {
		log.Println(<-receiver)
	}
}
