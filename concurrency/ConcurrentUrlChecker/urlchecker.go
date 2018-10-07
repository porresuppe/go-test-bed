// Check race conditions with go run -race urlchecker.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const checkInterval = 3
const statusInterval = 10

type urlState struct {
	url        string
	statusCode int
	latency    float64
}

func (u urlState) String() string {
	return fmt.Sprintf("Call to %s returned %d (took %f sec)", u.url, u.statusCode, u.latency)
}

func printState(stateChannel chan urlState) {
	state := make(map[string]urlState)
	start := time.Now()
	for {
		us := <-stateChannel
		state[us.url] = us

		delta := time.Since(start)
		if delta > statusInterval*time.Second {
			start = time.Now()
			fmt.Printf("%v\n", state)
		}
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s ran for %s", name, elapsed)
}

func checkURL(url string, stateChannel chan urlState) {
	start := time.Now()
	for {
		delta := time.Since(start)
		if delta < checkInterval*time.Second {
			continue
		}
		start = time.Now()

		log.Printf("Calling: %s\n", url)
		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			log.Fatalf("could not create request: %v", err)
		}
		client := http.DefaultClient

		statusCode := 503

		startOfCall := time.Now()
		res, err := client.Do(req)
		if err == nil {
			statusCode = res.StatusCode
		}
		elapsed := time.Since(startOfCall)
		stateChannel <- urlState{url, statusCode, elapsed.Seconds()}
	}
}

func main() {
	defer timeTrack(time.Now(), "urlchecker")
	urls := []string{"http://127.0.0.1:8080/hello/kalle", "http://127.0.0.1:8080/hello/anker"}

	stateChannel := make(chan urlState)

	go printState(stateChannel)

	for _, v := range urls {
		go checkURL(v, stateChannel)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}
		if char == 'q' {
			break
		}
	}
}
