package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/datacraft/deleggit/core/internal/domain"
)

// TrendScout is a "Communicator" agent (Level 2).
// It analyzes sensor data streams to find patterns and signal anomalies.
type TrendScout struct {
	id          string
	threshold   float64
	readings    []float64 // Moving window
	windowSize  int
}

func NewTrendScout(id string, threshold float64) *TrendScout {
	return &TrendScout{
		id:         id,
		threshold:  threshold,
		readings:   make([]float64, 0),
		windowSize: 5,
	}
}

func (a *TrendScout) ID() string {
	return a.id
}

func (a *TrendScout) Type() domain.AgentType {
	return domain.AgentTypeCommunicator
}

func (a *TrendScout) Execute(ctx context.Context, input domain.CloudEvent) (*domain.CloudEvent, error) {
	// 1. Parse Input (Expects: "sensor.cpu.temp")
	if input.Type != "sensor.cpu.temp" {
		return nil, nil // Ignore non-matching events
	}

	// Payload is e.g., "45.2 C" inside JSON string, or just raw bytes.
	// device-mock sends: `{"value": 45.2, "unit": "C"}` (Hypothetically)
	// Actually, looking at logs: `Data Length: 82`.
	// Let's assume the mock sends a simple JSON object.
	
	var data struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	}
	
	// device-mock currently might just be sending raw text or simple JSON.
	// Let's try to unmarshal.
	if err := json.Unmarshal(input.Data, &data); err != nil {
		// Fallback: If mock sends just a number string?
		str := string(input.Data) 
		// Strip quotes if it's a JSON string
		str = strings.Trim(str, "\"")
		val, err := strconv.ParseFloat(str, 64)
		if err == nil {
			data.Value = val
		} else {
			// Just log raw for debugging phase
			// log.Printf("[%s] Could not parse data: %s", a.ID(), string(input.Data))
			return nil, nil
		}
	}

	// 2. Logic: Moving Average & Threshold
	a.readings = append(a.readings, data.Value)
	if len(a.readings) > a.windowSize {
		a.readings = a.readings[1:]
	}

	// Calculate Avg
	sum := 0.0
	for _, v := range a.readings {
		sum += v
	}
	avg := sum / float64(len(a.readings))

	log.Printf("  >>> [AGENT:%s] Analysis: Current=%.1f Avg=%.1f (Threshold=%.1f)", a.id, data.Value, avg, a.threshold)

	// 3. Output: Signal if Overheat
	if avg > a.threshold {
		log.Printf("  !!! [AGENT:%s] OVERHEAT DETECTED. Signaling Swarm.", a.id)
		
		signalData := map[string]interface{}{
			"severity": "critical",
			"message":  fmt.Sprintf("CPU High Temp detected: %.1f C (Avg %.1f)", data.Value, avg),
			"source":   input.Source,
		}
		
		evt, err := domain.NewEvent("agent.trend_scout", "swarm.security.alert", signalData)
		if err != nil {
			return nil, err
		}
		return &evt, nil
	}

	// No anomaly
	return nil, nil
}
