package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func compareJSON(expected, actual string) bool {
	var expectedJSON, actualJSON interface{}

	err1 := json.Unmarshal([]byte(expected), &expectedJSON)
	err2 := json.Unmarshal([]byte(actual), &actualJSON)

	return err1 == nil && err2 == nil && expectedJSON == actualJSON
}

func TestStatusHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/app/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(statusHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"cpu":{"highPriority":0.68},"replicas":10}`
	if !compareJSON(expected, rr.Body.String()) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestReplicasHandler(t *testing.T) {
	var jsonStr = []byte(`{"replicas":15}`)
	req, err := http.NewRequest("PUT", "/app/replicas", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(replicasHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"message":"Replicas updated"}`
	if !compareJSON(expected, rr.Body.String()) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	
	if currentStatus.Replicas != 15 {
		t.Errorf("replicas not updated: got %v want %v",
			currentStatus.Replicas, 15)
	}
}
