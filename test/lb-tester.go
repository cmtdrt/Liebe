package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := "http://localhost:8080/users/1"
	total := 100
	concurrency := 10

	var wg sync.WaitGroup
	wg.Add(concurrency)

	requestsPerWorker := total / concurrency

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for w := 0; w < concurrency; w++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < requestsPerWorker; i++ {
				resp, err := client.Get(url)
				if err != nil {
					fmt.Printf("[worker %d] error: %v\n", id, err)
					continue
				}
				fmt.Printf("[worker %d] status: %d\n", id, resp.StatusCode)
				resp.Body.Close()
			}
		}(w)
	}

	wg.Wait()
	fmt.Println("Done.")
}
