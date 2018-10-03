package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

var shoppingList []shoppingItem

type shoppingItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Supermarket string  `json:"supermarket"`
	Price       float64 `json:"price"`
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		enc := json.NewEncoder(w)
		err := enc.Encode(shoppingList)
		if err != nil {
			// if encoding fails, create an error page with code 500.
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "POST":
		var item shoppingItem
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item.ID = fmt.Sprintf("%d", len(shoppingList))
		shoppingList = append(shoppingList, item)
	case "DELETE":
		re := regexp.MustCompile("^/items/([0-9]+)$")
		s := re.FindStringSubmatch(r.URL.Path)
		if s != nil {
			removeIndex := 0
			for i, v := range shoppingList {
				if s[1] == v.ID {
					removeIndex = i
					break
				}
			}
			shoppingList = append(shoppingList[:removeIndex], shoppingList[removeIndex+1:]...)
		} else {
			shoppingList = nil
		}
	default:
		http.Error(w, "StatusMethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}
}

func totalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		total := 0.0
		for _, v := range shoppingList {
			total += v.Price
		}
		enc := json.NewEncoder(w)
		err := enc.Encode(struct {
			Total float64 `json:"total price"`
		}{total})
		if err != nil {
			// if encoding fails, create an error page with code 500.
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "StatusMethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/items/", itemsHandler)
	mux.HandleFunc("/total", totalHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
