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
	jobs := make(chan int, 10000)
	results := make(chan int, 10000)

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
		for i := 1; i <= 100_000; i++ {
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
}

func executeRandom(snapshot func(i int)) {
	start := time.Now()

	// Number of goroutines to use (based on CPU cores)
	numWorkers := runtime.NumCPU()

	// Create channels for job distribution
	jobs := make(chan int, 10000)
	results := make(chan int, 10000)

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
		for i := 1; i <= 100_000; i++ {
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
}

func display(i int) {
}

func workerLoadTest(jobs <-chan int, wg *sync.WaitGroup, counter *uint64, requestFunc func()) {
	defer wg.Done()
	for range jobs {
		requestFunc() // Call function from main.go
		atomic.AddUint64(counter, 1)
	}
}

func loadTest(totalRequests int, numWorkers int, requestFunc func()) {
	start := time.Now()
	jobs := make(chan int, totalRequests)
	var wg sync.WaitGroup
	var counter uint64

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go workerLoadTest(jobs, &wg, &counter, requestFunc)
	}

	// Send jobs to workers
	for i := 0; i < totalRequests; i++ {
		jobs <- i
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()

	// Print performance results
	elapsed := time.Since(start)
	fmt.Printf("Processed %d requests in %s using %d workers\n",
		atomic.LoadUint64(&counter), elapsed, numWorkers)
	fmt.Printf("Achieved %.2f RPS\n",
		float64(atomic.LoadUint64(&counter))/elapsed.Seconds())
}
