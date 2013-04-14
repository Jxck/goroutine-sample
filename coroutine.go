package main

import (
	"log"
)

func f(yield chan string) {
	yield <- "one"
	yield <- "two"
	yield <- "three"
}

func main() {
	co := make(chan string)
	go f(co)
	log.Println(<-co) // one
	log.Println(<-co) // two
	log.Println(<-co) // three
}
