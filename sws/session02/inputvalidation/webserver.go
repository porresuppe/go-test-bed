package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		name = "friend"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func helloAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", http.StatusBadRequest)
		return
	}
	if r.Body == http.NoBody {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}
	var s string
	for key, value := range r.Form {
		fmt.Printf("Type: %T", value)
		s += fmt.Sprintf("%v: %v ", key, strings.Join(value, ", "))
	}

	fmt.Fprintf(w, "%s", s)
}

type textHandler struct {
	h http.HandlerFunc
}

func (t textHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set the content type
	w.Header().Set("Content-Type", "text/plain")
	// Then call ServeHTTP in the decorated handler.
	t.h(w, r)
}

type person struct {
	Name     string `json:"name"`
	AgeYears int    `json:"age_years"`
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	var p person

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Name is %v and age is %v", p.Name, p.AgeYears)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", decodeHandler)
	mux.Handle("/hello", textHandler{helloHandler})
	mux.HandleFunc("/helloall", helloAllHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
