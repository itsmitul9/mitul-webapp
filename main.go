package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"
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

const (
	targetCPUUsage  = 0.80 //target usgae set to 80%
	minReplicas     = 1
	maxReplicas     = 100
	monitorInterval = 2 * time.Second //auto-scaler check cpu usage at this interval
)

func main() {
	port := flag.String("port", "8123", "Port for the application")
	flag.Parse()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/app/status", statusHandler)
	http.HandleFunc("/app/replicas", replicasHandler)

	//start as go routine
	go autoScaler()

	log.Printf("Starting server on port %s\n", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	message := map[string]string{"message": "Welcome to the Auto-Scaler App !!!"}

	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// handles http GET request for /api/status
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	defer mutex.Unlock()
	json.NewEncoder(w).Encode(currentStatus)
}

// handles http GET & PUT for replicas update
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

		// check requested number of replicas
		if update.Replicas < minReplicas || update.Replicas > maxReplicas {
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

// core logic of code that check cpu usage & adjust the replica count to maintain desired cpu usage
func autoScaler() {
	for {
		time.Sleep(monitorInterval)

		mutex.Lock()

		currentStatus.CPU.HighPriority = simulateCPUUsage(currentStatus.Replicas)

		// added a buffer to avoid frequent scaling due to fluctuation in cpu usage
		targetBuffer := 0.02

		// adjust replicas when cpu usage increase or descrease
		if currentStatus.CPU.HighPriority > targetCPUUsage+targetBuffer {
			diff := currentStatus.CPU.HighPriority - targetCPUUsage
			// formula to calculate the auto-scaler function
			extrareplicas := int(math.Ceil(diff * float64(currentStatus.Replicas)))
			currentStatus.Replicas = currentStatus.Replicas + extrareplicas

			if currentStatus.Replicas > maxReplicas {
				currentStatus.Replicas = maxReplicas
			}
		} else if currentStatus.CPU.HighPriority < targetCPUUsage-targetBuffer && currentStatus.Replicas > minReplicas {
			diff := targetCPUUsage - currentStatus.CPU.HighPriority
			reducedreplicas := int(math.Ceil(diff * float64(currentStatus.Replicas)))
			currentStatus.Replicas = currentStatus.Replicas - reducedreplicas
			if currentStatus.Replicas < minReplicas {
				currentStatus.Replicas = minReplicas
			}
		}

		log.Printf("cpu usage:%.2f replicas:%d\n", currentStatus.CPU.HighPriority, currentStatus.Replicas)

		mutex.Unlock()
	}
}

func simulateCPUUsage(replicas int) float64 {
	// added flucuation with base cpu usage
	baseCPUUsage := 0.68 + rand.Float64()*0.1
	// logic to decrease cpu usage when replicas increase
	return baseCPUUsage + targetCPUUsage/float64(replicas)
}
