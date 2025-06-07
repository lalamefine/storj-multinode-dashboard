package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Constantes pour la configuration du cache
const (
	cacheExpiration   = 30 * time.Minute
	refreshInterval   = 30 * time.Minute
	backgroundRefresh = true // Active/désactive le rafraîchissement automatique
)

func main() {
	// Définir les gestionnaires de routes
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/health", healthCheckHandler)

	// Lancer la goroutine pour le rafraîchissement automatique du cache
	if backgroundRefresh {
		go autoRefreshCache()
	}

	// Configurer le port (8080 par défaut)
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	// Démarrer le serveur
	serverAddr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server listening on port %s\n", serverAddr)

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
