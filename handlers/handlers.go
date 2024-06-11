package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sakthe-Balan/GoMongoDB/db"
)

var database *db.Driver

func InitDB(dir string) {
	var err error
	database, err = db.New(dir, nil)
	if err != nil {
		fmt.Println("Error initializing database:", err)
	}
}

func CreateResourceHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := database.Write(collection, resource, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func ReadResourceHandler(w http.ResponseWriter, r *http.Request) {
	collection := r.URL.Query().Get("collection")
	resource := r.URL.Query().Get("resource")
	if collection == "" || resource == "" {
		http.Error(w, "Missing collection or resource name", http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	if err := database.Read(collection, resource, &data); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func ReadAllResourcesHandler(w http.ResponseWriter, r *http.Request) {
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		http.Error(w, "Missing collection name", http.StatusBadRequest)
		return
	}

	records, err := database.ReadAll(collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data []map[string]interface{}
	for _, record := range records {
		var item map[string]interface{}
		if err := json.Unmarshal(record, &item); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, item)
	}

	json.NewEncoder(w).Encode(data)
}

func DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	collection := r.URL.Query().Get("collection")
	resource := r.URL.Query().Get("resource")
	if collection == "" || resource == "" {
		http.Error(w, "Missing collection or resource name", http.StatusBadRequest)
		return
	}

	if err := database.Delete(collection, resource); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		http.Error(w, "Missing collection name", http.StatusBadRequest)
		return
	}

	if err := database.DeleteAll(collection); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var query map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := database.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func RegexSearchHandler(w http.ResponseWriter, r *http.Request) {
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		http.Error(w, "Missing collection name", http.StatusBadRequest)
		return
	}

	var query map[string]string
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := database.RegexSearch(collection, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}
