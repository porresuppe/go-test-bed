package main

// package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getName(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	fmt.Fprintln(w, "Hello, "+name)
}

func main() {
	r := mux.NewRouter()

	// match GET regardless of productID
	r.HandleFunc("/hello/{name}", getName)

	// handle all requests with the Gorilla router.
	http.Handle("/", r)
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal(err)
	}
}
