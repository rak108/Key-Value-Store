package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type key_value_store struct {
	database map[string]string
	mu       sync.RWMutex
}

func (key_value *key_value_store) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		key_value.mu.RLock()
		w.Write([]byte(`{"GET request successful"}`))
		key_value.mu.RUnlock()

	case "POST":
		w.WriteHeader(http.StatusCreated)
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//w.Write([]byte(`{"message": "post called"}`))
		key_value.mu.Lock()
		Key := r.FormValue("Key")
		Value := r.FormValue("Value")
		if _, check := key_value.database[Key]; !check {
			key_value.database[Key] = Value
			w.Write([]byte(`{"POST request successful: Created entry"}`))
		}
		key_value.mu.Unlock()

	case "PUT":
		w.WriteHeader(http.StatusAccepted)
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		key_value.mu.Lock()
		Key := r.FormValue("Key")
		Value := r.FormValue("Value")
		if _, check := key_value.database[Key]; check {
			key_value.database[Key] = Value
			w.Write([]byte(`{"PUT request successful: Created entry"}`))
		}
		else
		{
			w.Write([]byte(`{"Inavlid Key-Value Pair"}`))
		}
		key_value.mu.Unlock()

	case "DELETE":
		w.WriteHeader(http.StatusOK)
		key_value.mu.Lock()
		Key := r.URL.Path[8:]
		//Value := r.FormValue("Value")
		check:=delete(key_value.database,Key)
		if check==true{
			w.Write([]byte(`{"DELETE request successful: Created entry"}`))
		}
		else
		{
			w.Write([]byte(`{"Inavlid Key-Value Pair"}`))
		}
		key_value.mu.Unlock()

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}

}

func main() {

	key_value = &key_value_store{}
    http.HandleFunc("/",key_value.ServeHTTP)
	//Start the server and listen for requests
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
