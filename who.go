package main

import (
	"fmt"
	"runtime"
)

func main() {
	buf := make([]byte, 1<<20)
	buf = buf[:runtime.Stack(buf, true)]
	fmt.Println(string(buf))
	select {}
}
