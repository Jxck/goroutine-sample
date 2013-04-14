package main

import (
	"log"
)

func generator(n int) chan int {
	ch := make(chan int)
	i := 0
	go func() {
		for {
			ch <- i
			i++
			if i > n {
				close(ch)
				break
			}
		}
	}()
	return ch
}

func main() {
	for x := range generator(10) {
		log.Println(x)
	}
}
