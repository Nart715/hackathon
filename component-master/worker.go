package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var accountIds = []int64{}

// Worker pool for parallel processing
func worker(jobs <-chan int, results chan<- int, wg *sync.WaitGroup, snapshot func(i int)) {
	defer wg.Done()
	for n := range jobs {
		fmt.Println("NUMBER >> ", n)
		snapshot(n)
		results <- n
	}
}

func execute(snapshot func(i int)) {
	start := time.Now()

	// Number of goroutines to use (based on CPU cores)
	numWorkers := runtime.NumCPU()

	// Create channels for job distribution
	jobs := make(chan int, 1000)
	results := make(chan int, 1000)

	// Create a wait group to track workers
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg, snapshot)
	}

	// Counter for tracking progress
	var counter uint64

	// Start a goroutine to send jobs
	go func() {
		for i := 1; i <= 10_000; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Start a goroutine to collect results and update counter
	go func() {
		for n := range results {
			accountIds = append(accountIds, int64(n))
			atomic.AddUint64(&counter, 1)
		}
	}()

	// Wait for all workers to complete
	wg.Wait()
	close(results)

	// Print results
	elapsed := time.Since(start)
	fmt.Printf("Processed %d numbers in %s using %d workers\n",
		atomic.LoadUint64(&counter), elapsed, numWorkers)
	fmt.Printf("Average processing time per number: %v\n",
		elapsed/time.Duration(atomic.LoadUint64(&counter)))
	fmt.Println("AccountIds: ", len(accountIds))
}

func display(i int) {
}
