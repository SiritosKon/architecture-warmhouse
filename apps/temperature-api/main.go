package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func resolve(location, sensorID string) (string, string) {
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}
	return location, sensorID
}

func newReading(location, sensorID string) TemperatureResponse {
	location, sensorID = resolve(location, sensorID)
	value := 18.0 + rand.Float64()*8.0
	value = float64(int(value*10)) / 10
	return TemperatureResponse{
		Value:       value,
		Unit:        "celsius",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: fmt.Sprintf("Temperature reading for %s", location),
	}
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	sensorID := r.URL.Query().Get("sensorId")
	writeJSON(w, newReading(location, sensorID))
}

func handlePath(w http.ResponseWriter, r *http.Request) {
	sensorID := strings.TrimPrefix(r.URL.Path, "/temperature/")
	if sensorID == "" {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, newReading("", sensorID))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/temperature", handleQuery)
	mux.HandleFunc("/temperature/", handlePath)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	addr := ":" + port
	log.Printf("temperature-api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
