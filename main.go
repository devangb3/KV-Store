package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("recieved request from: %v", r.URL.Path)
	fmt.Fprintf(w, "Printing %s", r.URL.Path[1:])
}

var kvStore = make(map[string]string)

func putKvStoreHandler(w http.ResponseWriter, r *http.Request) {
	var data struct{
		Key string `json:"key"`
		Value string `json:"value`
	}
	
	err := json.NewDecoder(r.Body).Decode(&data);
	if err != nil{
		fmt.Fprintf(w, "Error in decoding!")
		return
	}
	if data.Key == "" || data.Value == ""{
		fmt.Fprintf(w, "Both key and value should be non-empty");
	}
	kvStore[data.Key] = data.Value;
	fmt.Fprintf(w, "Key : %v stored with value : %v", data.Key, data.Value)
}

func getKVStoreHandler(w http.ResponseWriter, r *http.Request) {
	
	key := r.URL.Query().Get("key")
	if key == ""{
		http.Error(w, "Key parameter empty", http.StatusBadRequest);
	}
	
	val, exists := kvStore[key];
	if !exists{
		http.Error(w,"{No record attached with key found}", http.StatusBadRequest);
		return;
	}
	response := struct {
		Value string
	}{
		Value: val,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response);
}
func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/put", putKvStoreHandler)
	http.HandleFunc("/get", getKVStoreHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
