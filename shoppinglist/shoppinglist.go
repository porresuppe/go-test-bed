package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type shoppingListItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Supermarket string  `json:"supermarket"`
	Price       float64 `json:"price"`
}

var shoppingList []shoppingListItem

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(w)
		err := enc.Encode(shoppingList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "POST":
		var item shoppingListItem

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id := 0
		if shoppingList != nil {
			id = shoppingList[len(shoppingList)-1].ID + 1
		}

		item.ID = id
		shoppingList = append(shoppingList, item)
	case "DELETE":
		shoppingList = nil
	default:
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
	}
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := string(vars["id"])

	for i, v := range shoppingList {
		if strconv.Itoa(v.ID) == id {
			shoppingList = append(shoppingList[:i], shoppingList[i+1:]...)
			break
		}
	}
}

func totalPriceHandler(w http.ResponseWriter, r *http.Request) {
	totalPrice := 0.0
	for _, v := range shoppingList {
		totalPrice += v.Price
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	err := enc.Encode(struct {
		TotalPrice float64 `json:"total price"`
	}{totalPrice})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func singleSupermarketListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	supermarket := string(vars["supermarket"])

	var singleSupermarketList []shoppingListItem
	for _, v := range shoppingList {
		if v.Supermarket == supermarket {
			singleSupermarketList = append(singleSupermarketList, v)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	err := enc.Encode(singleSupermarketList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/items", itemsHandler).Methods("GET", "POST", "DELETE")
	r.HandleFunc("/items/{id:[0-9]+}", itemHandler).Methods("DELETE")
	r.HandleFunc("/items/totprice", totalPriceHandler).Methods("GET")
	r.HandleFunc("/items/{supermarket}", singleSupermarketListHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
