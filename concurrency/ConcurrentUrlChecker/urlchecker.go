// Check race conditions with go run -race urlchecker.go
// TODO: Remove mutex - see https://gobyexample.com/stateful-goroutines
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const checkInterval = 3
const statusInterval = 10

var state map[string]urlState

type urlState struct {
	statusCode int
	latency    float64
}

func (u urlState) String() string {
	return fmt.Sprintf(" returned %d (took %f sec)", u.statusCode, u.latency)
}

var mutex = &sync.Mutex{}

func main() {
	urls := []string{"http://127.0.0.1:8080/hello/kalle", "http://127.0.0.1:8080/hello/anker"}
	urlcheckers := 2
	state = make(map[string]urlState)
	ch := make(chan string)

	go printState()

	for {
		for _, v := range urls {
			go func(url string) {
				ch <- url
			}(v)
		}

		for i := 1; i < urlcheckers; i++ {
			go checkURL(ch)
		}

		time.Sleep(1 * time.Second) // Slow down main loop
	}
}

func printState() {
	for {
		time.Sleep(statusInterval * time.Second)
		mutex.Lock()
		fmt.Printf("%v\n", state)
		mutex.Unlock()
	}
}

func checkURL(ch chan string) {
	time.Sleep(checkInterval * time.Second)
	url := <-ch
	log.Printf("Calling: %s\n", url)
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Fatalf("could not create request: %v", err)
	}
	client := http.DefaultClient

	statusCode := 503

	start := time.Now()
	res, err := client.Do(req)
	if err == nil {
		statusCode = res.StatusCode
	}
	elapsed := time.Since(start)
	mutex.Lock()
	state[url] = urlState{statusCode, elapsed.Seconds()}
	mutex.Unlock()
	log.Printf("Done calling %s\n", url)
}
