package main

import (
	"fmt"
	"net/http"

	"github.com/Sakthe-Balan/GoMongoDB/handlers"
)

func main() {
	dir := "./dbase"
	handlers.InitDB(dir)

	http.HandleFunc("/write", handlers.CreateResourceHandler)     // POST
	http.HandleFunc("/read", handlers.ReadResourceHandler)        // GET
	http.HandleFunc("/readall", handlers.ReadAllResourcesHandler) // GET
	http.HandleFunc("/delete", handlers.DeleteResourceHandler)    // DELETE
	http.HandleFunc("/deleteall", handlers.DeleteAllHandler)      // DELETE
	http.HandleFunc("/search", handlers.SearchHandler)            // POST

	fmt.Println("Starting server on :6942")
	if err := http.ListenAndServe(":6942", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
