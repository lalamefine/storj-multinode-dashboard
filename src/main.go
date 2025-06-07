package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	cacheExpiration   = 30 * time.Minute
	refreshInterval   = 30 * time.Minute
	backgroundRefresh = true
)

func main() {
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/health", healthCheckHandler)

	if backgroundRefresh {
		go autoRefreshCache()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	serverAddr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("Server listening on %s\n", serverAddr)

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
