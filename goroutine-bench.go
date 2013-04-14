package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	count := 1000 * 100

	var startMemory runtime.MemStats
	runtime.ReadMemStats(&startMemory)

	start := time.Now()
	fin := make(chan bool)
	for i := 0; i < count; i++ {
		go func() {
			<-fin
		}()
	}
	elapsed := time.Since(start)

	var endMemory runtime.MemStats
	runtime.ReadMemStats(&endMemory)
	close(fin)

	fmt.Printf(`
goroutine:  %d
cpu:        %d
time:       %f
memory all: %f MB
memory:     %f byte/goroutine
`,
		count, cpus, elapsed.Seconds(),
		float64(endMemory.Alloc-startMemory.Alloc)/float64(1024*1024),
		float64(endMemory.Alloc-startMemory.Alloc)/float64(count))
}
