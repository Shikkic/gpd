package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	recurTime := recur()
	concurrentRecurTime := concurrentRecur()
	log.Println("")
	log.Printf("Recur finished in %f", recurTime.Seconds())
	log.Printf("Concurrent Recur finished in %f", concurrentRecurTime.Seconds())
}

func recur() time.Duration {
	startTime := time.Now()
	dir, _ := os.Getwd()
	recursiveReadDir(dir)
	return time.Since(startTime)
}

func recursiveReadDir(dir string) {
	log.Printf("Recur:\t%s\n", dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			newDir := fmt.Sprintf("%s/%s", dir, file.Name())
			recursiveReadDir(newDir)
		}
	}
}

func concurrentRecur() time.Duration {
	startTime := time.Now()
	// This will tell us when we're done.
	var wg sync.WaitGroup

	dir, _ := os.Getwd()
	wg.Add(1)
	concurrentRecursiveReadDir(dir, &wg)

	wg.Wait()
	return time.Since(startTime)
}

func concurrentRecursiveReadDir(dir string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("ConcurrentRecur:\t%s\n", dir)
	files, err := ioutil.ReadDir(dir)

	// Panic if we have any unexpected errors.
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			newDir := fmt.Sprintf("%s/%s", dir, file.Name())
			wg.Add(1)
			go concurrentRecursiveReadDir(newDir, wg)
		}
	}
}
