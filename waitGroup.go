package main

import (
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			log.Println(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
