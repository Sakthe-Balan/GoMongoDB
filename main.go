package main

import (
	"fmt"
	"net/http"

	"github.com/Sakthe-Balan/GoMongoDB/handlers"
)

func main() {
	dir := "./dbase"
	handlers.InitDB(dir)

	http.HandleFunc("/write", handlers.CreateResourceHandler)
	http.HandleFunc("/read", handlers.ReadResourceHandler)
	http.HandleFunc("/readall", handlers.ReadAllResourcesHandler)
	http.HandleFunc("/delete", handlers.DeleteResourceHandler)

	fmt.Println("Starting server on :6942")
	if err := http.ListenAndServe(":6942", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
