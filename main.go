package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/devangb3/KV-Store/config"
	"github.com/devangb3/KV-Store/database"
	"github.com/joho/godotenv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("recieved request from: %v", r.URL.Path)
	fmt.Fprintf(w, "Printing %s", r.URL.Path[1:])
}

type Server struct{
	store *database.Store
}
func NewServer(store *database.Store) *Server{
	return &Server{store:store};
}

func(s *Server) putKvStoreHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := s.store.InsertRecord(data.Key, data.Value); err!=nil{
		log.Fatalf("Error while inserting record : %v\n", err);
		http.Error(w, "Could not insert record", http.StatusBadRequest);
		return;
	}

	log.Printf("Key : %v stored with value : %v", data.Key, data.Value)
	
	response := struct{
		Value string
	}{
		Value : fmt.Sprintf("Stored key %v successfully with value %v", data.Key, data.Value),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response);
}

func(s *Server) getKVStoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet{
		http.Error(w,"Invalid method, only accepts GET requests", http.StatusMethodNotAllowed);
		return;
	}
	key := r.URL.Query().Get("key")
	if key == ""{
		http.Error(w, "Key parameter empty", http.StatusBadRequest);
	}
	val,err := s.store.GetRecord(key);
	if err != nil{
		log.Fatalf("Error while getting key : %v", err)
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

func(s *Server) deleteKVStoreHandler(w http.ResponseWriter, r *http.Request){
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
	log.Printf("Recieved delete request for key : %v\n", data.Key);

	if err := s.store.DeleteRecord(data.Key); err !=nil{
		log.Fatalf("Error while deleting record: %v\n", err)
		http.Error(w, "Error while deleting record", http.StatusBadRequest);
		return;
	}

	log.Printf("Key Deleted successfully %v\n", data.Key);
	response := struct {
		Data string
	}{
		Data: "Deleted Successfully!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response);
}


func main() {
	
	err := godotenv.Load();
	if err != nil{
		log.Fatalf("error Loading env variables %v\n", err);
		return;
	}
	cfg, err := config.LoadConfig();
	if err != nil{
		log.Fatalf("Error loading config %v\n", err)
		return;
	}
	store,err := database.NewStore(*cfg);
	if err != nil{
		log.Fatalf("Error Creating New Store %v\n", err);
		return;
	}
	defer store.Close();

	server := NewServer(store);

	http.HandleFunc("/", handler)
	http.HandleFunc("/put", server.putKvStoreHandler)
	http.HandleFunc("/get", server.getKVStoreHandler)
	http.HandleFunc("/delete", server.deleteKVStoreHandler)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
