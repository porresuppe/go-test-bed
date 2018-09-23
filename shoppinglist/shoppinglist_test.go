package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestItemsHandler(t *testing.T) {
	{
		req, err := http.NewRequest("POST", "/items", strings.NewReader(`{"name":"beer","supermarket":"netto","price": 8}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/json")

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(itemsHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}
	{
		req, err := http.NewRequest("GET", "/items", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(itemsHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := `[{"id":0,"name":"beer","supermarket":"netto","price":8}]
` // actual has a newline character
		actual := rr.Body.String()
		if actual != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				actual, expected)
		}
	}
}
