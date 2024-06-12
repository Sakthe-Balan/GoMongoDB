package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sakthe-Balan/GoMongoDB/Distributed_Framework/db"
)

var distributedDatabase *db.DistributedDriver

func InitDistributedDB(nodes []string) {
	var err error
	distributedDatabase, err = db.NewDistributedDriver(nodes)
	if err != nil {
		fmt.Println("Error initializing distributed database:", err)
	}
}

func CreateDistributedResourceHandler(w http.ResponseWriter, r *http.Request) {
	collection := r.URL.Query().Get("collection")
	resource := r.URL.Query().Get("resource")
	if collection == "" || resource == "" {
		http.Error(w, "Missing collection or resource name", http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := distributedDatabase.Write(collection, resource, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func ReadDistributedResourceHandler(w http.ResponseWriter, r *http.Request) {
	collection := r.URL.Query().Get("collection")
	resource := r.URL.Query().Get("resource")
	if collection == "" || resource == "" {
		http.Error(w, "Missing collection or resource name", http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	if err := distributedDatabase.Read(collection, resource, &data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func ReadAllDistributedResourcesHandler(w http.ResponseWriter, r *http.Request) {
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		http.Error(w, "Missing collection name", http.StatusBadRequest)
		return
	}

	data, err := distributedDatabase.ReadAll(collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func DeleteDistributedResourceHandler(w http.ResponseWriter, r *http.Request) {
	collection := r.URL.Query().Get("collection")
	resource := r.URL.Query().Get("resource")
	if collection == "" || resource == "" {
		http.Error(w, "Missing collection or resource name", http.StatusBadRequest)
		return
	}

	if err := distributedDatabase.Delete(collection, resource); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SearchDistributedResourcesHandler(w http.ResponseWriter, r *http.Request) {
	var query map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := distributedDatabase.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}
