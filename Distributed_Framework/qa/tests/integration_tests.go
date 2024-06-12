package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sakthe-Balan/GoMongoDB/Distributed_Framework/handlers"
)

func TestDistributedWrite(t *testing.T) {
	data := map[string]interface{}{
		"name":  "Test User",
		"email": "test@example.com",
	}
	body, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "/distributed/write?collection=users&resource=testuser", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateDistributedResourceHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestDistributedRead(t *testing.T) {
	req, err := http.NewRequest("GET", "/distributed/read?collection=users&resource=testuser", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ReadDistributedResourceHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if response["name"] != "Test User" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response["name"], "Test User")
	}
}

// Add more tests for read all, delete, and search handlers
