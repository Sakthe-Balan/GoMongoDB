package main

import (
	"fmt"
	"net/http"

	"github.com/Sakthe-Balan/GoMongoDB/handlers"
)

func main() {
	dir := "./dbase"
	handlers.InitDB(dir)

	http.HandleFunc("/create", handlers.CreateUserHandler)
	http.HandleFunc("/read", handlers.ReadUserHandler)
	http.HandleFunc("/readall", handlers.ReadAllUsersHandler)
	http.HandleFunc("/delete", handlers.DeleteUserHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
