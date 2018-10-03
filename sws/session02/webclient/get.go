package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	// try changing the value of this url
	res, err := http.Get("https://golang.org")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		fmt.Println(res.Status)
	} else {

		// TODO: read into dyn array
		// b := make([]byte, 256)
		// for {
		// 	n, err := res.Body.Read(b)
		// 	if n == len(b) {

		// 	}
		// 	if err == io.EOF {
		// 		break
		// 	}
		// }

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", b)

	}
	doPut()
	doGet()
}

func doPut() {
	req, err := http.NewRequest("PUT", "https://http-methods.appspot.com/lirumlarum/Message", strings.NewReader("Kalle"))
	if err != nil {
		log.Fatalf("could not create request: %v", err)
	}
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("http request failed: %v", err)
	}
	fmt.Println(res.Status)
}

func doGet() {
	req, err := http.NewRequest("GET", "https://http-methods.appspot.com/lirumlarum/", nil)
	if err != nil {
		log.Fatalf("could not create request: %v", err)
	}
	req.URL.RawQuery = "v=true"
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("http request failed: %v", err)
	}
	defer res.Body.Close()
	fmt.Println(res.Status)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", b)
}
