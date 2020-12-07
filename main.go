package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type key_value_store struct {
	database map[string]string
	mu       sync.RWMutex
}

func (key_value *key_value_store) getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	key_value.mu.RLock()
	//fmt.Fprintf(w, "GET request accepted:\n")
	identify := mux.Vars(r)
	key := identify["key"]
	//fmt.Fprintf(w, "For Key: %s,\n", key)
	if _, check := key_value.database[key]; check {
		value := key_value.database[key]
		//fmt.Fprintf(w, "Value: %s\n", value)
		w.Write([]byte(fmt.Sprintf(`{"%s": "%s"}`, key, value)))
	} else {
		//fmt.Fprintf(w, "Key-Value pair does not exist\n")
		w.Write([]byte(`{"Error": "Key-Value pair not found"}`))
	}
	key_value.mu.RUnlock()
}

func (key_value *key_value_store) postHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	key_value.mu.Lock()
	//fmt.Fprintf(w, "POST request accepted:\n")
	identify := mux.Vars(r)
	key := identify["key"]
	value := identify["value"]
	//fmt.Fprintf(w, "For Key: %s,\n", key)
	//fmt.Fprintf(w, "Value: %s\n", value)
	if _, check := key_value.database[key]; !check {
		key_value.database[key] = value
		//fmt.Fprintf(w, "Key-Value store updated.\n")
		w.Write([]byte(fmt.Sprintf(`{"%s": "%s"}`, key, value)))
	} else {
		//fmt.Fprintf(w, "Key-Value pair already exists.\n")
		w.Write([]byte(`{"Error": "Key-Value pair already exists."}`))
	}
	key_value.mu.Unlock()
}

func (key_value *key_value_store) putHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	key_value.mu.Lock()
	//fmt.Fprintf(w, "PUT request accepted:\n")
	identify := mux.Vars(r)
	key := identify["key"]
	value := identify["value"]
	//fmt.Fprintf(w, "For Key: %s,\n", key)
	if _, check := key_value.database[key]; check {
		key_value.database[key] = value
		//fmt.Fprintf(w, "Value: %s\n", value)
		//fmt.Fprintf(w, "Key-Value store updated.\n")
		//w.Write([]byte(`{"Error": "Key-Value pair not found"}`))
		w.Write([]byte(fmt.Sprintf(`{"%s": "%s"}`, key, value)))
	} else {
		//fmt.Fprintf(w, "Key-Value pair does not exist.\n")
		w.Write([]byte(`{"Error": "Key-Value pair not found"}`))
	}
	key_value.mu.Unlock()
}

func (key_value *key_value_store) delHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	key_value.mu.Lock()
	//fmt.Fprintf(w, "DELETE request accepted:\n")
	identify := mux.Vars(r)
	key := identify["key"]
	//fmt.Fprintf(w, "Key: %s\n", key)
	if _, check := key_value.database[key]; check {
		delete(key_value.database, key)
		//	fmt.Fprintf(w, "Key-Value pair deleted.")
		w.Write([]byte(`{"Message": "Key-Value pair deleted."}`))
	} else {
		w.Write([]byte(`{"Error": "Key-Value pair not found"}`))
	}
	key_value.mu.Unlock()
}

func synerror(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"Message": "Not Found"}`))
}

func (key_value *key_value_store) getExisting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//fmt.Fprintf(w, "Key-Value pairs:\n")
	json.NewEncoder(w).Encode(key_value.database)
}

func main() {

	key_value := &key_value_store{}
	key_value.database = make(map[string]string)

	r := mux.NewRouter()
	r.HandleFunc("/existing", key_value.getExisting).Methods("GET")
	r.HandleFunc("/{key}", key_value.getHandler).Methods(http.MethodGet)
	r.HandleFunc("/{key}-{value}", key_value.postHandler).Methods("POST")
	r.HandleFunc("/{key}-{value}", key_value.putHandler).Methods(http.MethodPut)
	r.HandleFunc("/{key}", key_value.delHandler).Methods(http.MethodDelete)
	r.HandleFunc("/", synerror)

	fmt.Printf("Starting server at port 8080...\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
