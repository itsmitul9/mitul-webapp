package main

import (
	"testing"
)

func TestAutoScaler(t *testing.T) {
	tests := []struct {
		initialReplicas  int
		cpuUsage         float64
		expectedReplicas int
	}{
		{10, 0.90, 15}, // Expect an increase in replicas
		{10, 0.70, 7},  // Expect a decrease in replicas
		{10, 0.80, 10}, // No change expected
	}

	for _, test := range tests {
		currentStatus.Replicas = test.initialReplicas
		currentStatus.CPU.HighPriority = test.cpuUsage

		autoScaler()

		if currentStatus.Replicas != test.expectedReplicas {
			t.Errorf("Expected %d replicas, but got %d", test.expectedReplicas, currentStatus.Replicas)
		}
	}
}

func TestSimulateCPUUsage(t *testing.T) {
	replicas := 10
	usage := simulateCPUUsage(replicas)

	if usage < 0 || usage > 1 {
		t.Errorf("CPU usage should be between 0 and 1, got %f", usage)
	}
}
