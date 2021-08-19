package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HealthCheck struct {
	Status string `json:"status"`
}

func init() {
	stars = make(map[string]int64)
}

func main() {
	http.HandleFunc("/landscape", landscapeHandler)
	http.HandleFunc("/health", healthCheckHandler)
	http.Handle("/", http.FileServer(http.Dir("./web/static")))

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func landscapeHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := getCicdProjects()
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(HealthCheck{
			Status: "Error",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
	return
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(HealthCheck{
		Status: "Healthy",
	})
	return
}
