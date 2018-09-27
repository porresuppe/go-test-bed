package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type shoppingListItemPost struct {
	Name        string  `json:"name"`
	Supermarket string  `json:"supermarket"`
	Price       float64 `json:"price"`
	Content     string  `json:"content"`
}

type shoppingListItem struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Supermarket string  `json:"supermarket"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageurl"`
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	switch r.Method {
	case "GET":
		q := datastore.NewQuery("shoppingListItem")
		var shoppingList []shoppingListItem
		_, err := q.GetAll(ctx, &shoppingList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(w)
		err = enc.Encode(shoppingList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "POST":
		var pItem shoppingListItemPost

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&pItem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		item := shoppingListItem{Name: pItem.Name, Supermarket: pItem.Supermarket, Price: pItem.Price}

		// https://cloud.google.com/appengine/docs/standard/go/datastore/entities#Go_Assigning_identifiers
		l, _, err := datastore.AllocateIDs(ctx, "shoppingListItem", nil, 1)
		key := datastore.NewKey(ctx, "shoppingListItem", "", l, nil)
		item.ID = key.IntID()
		item.ImageURL, err = uploadFile(ctx, &pItem, &item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		key, err = datastore.Put(ctx, key, &item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "DELETE":
		q := datastore.NewQuery("shoppingListItem")
		var shoppingList []shoppingListItem
		keys, err := q.GetAll(ctx, &shoppingList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = datastore.DeleteMulti(ctx, keys)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
	}
}

func uploadFile(ctx context.Context, pItem *shoppingListItemPost, item *shoppingListItem) (url string, err error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	bucket := "ShoppingListImageBucket"
	name := fmt.Sprintf("%d-%s", item.ID, item.Name)
	wc := client.Bucket(bucket).Object(name).NewWriter(ctx)

	// Warning: storage.AllUsers gives public read access to anyone.
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	wc.ContentType = "image/jpg"

	// Entries are immutable, be aggressive about caching (1 day).
	wc.CacheControl = "public, max-age=86400"

	b, err := base64.StdEncoding.DecodeString(pItem.Content)
	if err != nil {
		return "", err
	}

	if _, err = io.Copy(wc, bytes.NewReader(b)); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(publicURL, bucket, name), nil
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	vars := mux.Vars(r)
	sID := vars["id"]
	id, err := strconv.ParseInt(sID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key := datastore.NewKey(ctx, "shoppingListItem", "", id, nil)
	err = datastore.Delete(ctx, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func totalPriceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	q := datastore.NewQuery("shoppingListItem")
	var shoppingList []shoppingListItem
	_, err := q.GetAll(ctx, &shoppingList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPrice := 0.0
	for _, v := range shoppingList {
		totalPrice += v.Price
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	err = enc.Encode(struct {
		TotalPrice float64 `json:"total price"`
	}{totalPrice})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func singleSupermarketListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	supermarket := string(vars["supermarket"])

	ctx := appengine.NewContext(r)

	q := datastore.
		NewQuery("shoppingListItem").
		Filter("Supermarket = ", supermarket)
	var shoppingList []shoppingListItem
	_, err := q.GetAll(ctx, &shoppingList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	err = enc.Encode(shoppingList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/items", itemsHandler).Methods("GET", "POST", "DELETE")
	r.HandleFunc("/items/{id:[0-9]+}", itemHandler).Methods("DELETE")
	r.HandleFunc("/items/totprice", totalPriceHandler).Methods("GET")
	r.HandleFunc("/items/{supermarket}", singleSupermarketListHandler).Methods("GET")

	http.Handle("/", r) // https://stackoverflow.com/a/26581607
}
