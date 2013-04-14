package main

import (
	"log"
)

func main() {
	fin := make(chan bool)
	go func() {
		log.Println("worker working..")
		close(fin) // fin <- false
	}()
	v, ok := <-fin
	log.Println(v)
	log.Println(ok)
}
