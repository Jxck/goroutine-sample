package main

import (
	"log"
	"runtime"
	"time"
)

func main() {
	go func() {
		log.Println("end")
	}()
	go func() {
		log.Println("return")
		return
	}()
	go func() {
		log.Println("exit")
		runtime.Goexit()
	}()
	time.Sleep(time.Second)
}
