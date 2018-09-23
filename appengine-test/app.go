package app

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// first create a new context
	c := appengine.NewContext(r)
	// and use that context to create a new http client
	client := urlfetch.Client(c)

	// now we can use that http client as before
	res, err := client.Get("http://google.com")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get google: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Got Google with status %s\n", res.Status)
}