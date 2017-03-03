package main

import (
	"fmt"

	"github.com/jsam/simplepool"
)

func printNumberJob(args ...interface{}) {
	fmt.Printf("%d\n", args...)
}

func main() {
	pool := simplepool.NewPool(100, 50)
	defer pool.Release()

	pool.WaitCount(10)

	for i := 0; i < 10; i++ {
		num := i
		pool.Enqueue(printNumberJob, num)
	}

	pool.WaitAll()
}
