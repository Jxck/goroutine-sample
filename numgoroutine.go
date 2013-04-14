package main

import (
	"log"
	"runtime"
)

func main() {
	log.Println(runtime.NumGoroutine())
	buf := make([]byte, 1<<20)
	buf = buf[:runtime.Stack(buf, true)]
	log.Println(string(buf))
	select {}
}
