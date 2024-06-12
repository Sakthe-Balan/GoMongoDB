package main

import (
	"fmt"
	"net/http"

	"github.com/Sakthe-Balan/GoMongoDB/Distributed_Framework/handlers"
	"github.com/Sakthe-Balan/GoMongoDB/Distributed_Framework/qa/monitoring"
)

func main() {
	nodes := []string{"node1:7001", "node2:7002", "node3:7003"}
	handlers.InitDistributedDB(nodes)

	http.HandleFunc("/distributed/write", handlers.CreateDistributedResourceHandler)
	http.HandleFunc("/distributed/read", handlers.ReadDistributedResourceHandler)
	http.HandleFunc("/distributed/readall", handlers.ReadAllDistributedResourcesHandler)
	http.HandleFunc("/distributed/delete", handlers.DeleteDistributedResourceHandler)
	http.HandleFunc("/distributed/search", handlers.SearchDistributedResourcesHandler)

	go monitoring.MonitorNodes()

	fmt.Println("Starting server on :6942")
	if err := http.ListenAndServe(":6942", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
