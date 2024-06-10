package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sakthe-Balan/GoMongoDB/db"
	"github.com/Sakthe-Balan/GoMongoDB/models"
)

var database *db.Driver

func InitDB(dir string) {
	var err error
	database, err = db.New(dir, nil)
	if err != nil {
		fmt.Println("Error initializing database:", err)
	}
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.Write("users", user.Name, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func ReadUserHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing user name", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.Read("users", name, &user); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func ReadAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	records, err := database.ReadAll("users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var users []models.User
	for _, record := range records {
		var user models.User
		if err := json.Unmarshal([]byte(record), &user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing user name", http.StatusBadRequest)
		return
	}

	if err := database.Delete("users", name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
