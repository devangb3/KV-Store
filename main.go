package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"sync"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("recieved request from: %v", r.URL.Path)
	fmt.Fprintf(w, "Printing %s", r.URL.Path[1:])
}

var (
	kvStore = make(map[string]string)
	mu sync.RWMutex
)

func putKvStoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost{
		http.Error(w, "Invalid method, only accepts POST requests", http.StatusMethodNotAllowed);
		return;
	}
	var data struct{
		Key string `json:"key"`
		Value string `json:"value"`
	}
	
	err := json.NewDecoder(r.Body).Decode(&data);
	if err != nil{
		http.Error(w, "Error occurred during Decoding", http.StatusBadRequest);
		return
	}
	if data.Key == "" || data.Value == ""{
		http.Error(w, "Either key or value is empty", http.StatusBadRequest);
		return;
	}
	mu.Lock()
	kvStore[data.Key] = data.Value;
	mu.Unlock();

	log.Printf("Key : %v stored with value : %v", data.Key, data.Value)
	
	response := struct{
		Value string
	}{
		Value : fmt.Sprintf("Stored key %v successfully with value %v", data.Key, data.Value),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response);
}

func getKVStoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet{
		http.Error(w,"Invalid method, only accepts GET requests", http.StatusMethodNotAllowed);
		return;
	}
	key := r.URL.Query().Get("key")
	if key == ""{
		http.Error(w, "Key parameter empty", http.StatusBadRequest);
	}
	mu.RLock();
	val,ok := kvStore[key];
	mu.RUnlock();
	if !ok{
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

func deleteKVStoreHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodDelete{
		http.Error(w, "Invalid Request method", http.StatusMethodNotAllowed)
		return;
	}
	var data struct{
		Key string `json:"key"`
	}
	err := json.NewDecoder(r.Body).Decode(&data);
	if err != nil{
		http.Error(w, "Error decoding JSON body", http.StatusBadRequest);
		return;
	}
	if(data.Key == ""){
		http.Error(w, "Key cannot be empty", http.StatusBadRequest);
		return;
	}
	log.Printf("Recieved delete request for key : ", data.Key);

	mu.Lock()
	delete(kvStore, data.Key);
	mu.Unlock()

	log.Printf("Key Deleted successfully %v", data.Key);
	response := struct {
		Data string
	}{
		Data: "Deleted Successfully!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response);
}


func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/put", putKvStoreHandler)
	http.HandleFunc("/get", getKVStoreHandler)
	http.HandleFunc("/delete", deleteKVStoreHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
