package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
)

type AppStatus struct {
	CPU      CPUStatus `json:"cpu"`
	Replicas int       `json:"replicas"`
}

type CPUStatus struct {
	HighPriority float64 `json:"highPriority"`
}

type ReplicaUpdate struct {
	Replicas int `json:"replicas"`
}

var (
	currentStatus = AppStatus{
		CPU:      CPUStatus{HighPriority: 0.68},
		Replicas: 10,
	}
	mutex sync.Mutex
)

func main() {
	port := flag.String("port", "8123", "Port for the application")
	flag.Parse()

	http.HandleFunc("/app/status", statusHandler)
	http.HandleFunc("/app/replicas", replicasHandler)

	log.Printf("Starting server on port %s\n", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	defer mutex.Unlock()
	json.NewEncoder(w).Encode(currentStatus)
}

func replicasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"replicas": currentStatus.Replicas})

	case http.MethodPut:
		var update ReplicaUpdate
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if update.Replicas < 1 {
			http.Error(w, "Invalid replicas count", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		currentStatus.Replicas = update.Replicas
		mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Replicas updated"})
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
