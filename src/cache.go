package main

import (
	"log"
	"sync"
	"time"
)

// Structure pour stocker les dernières requêtes
type lastRequest struct {
	names     []string
	endpoints []string
	timestamp time.Time
	mutex     sync.RWMutex
}

// Structure de cache pour stocker les données des nœuds
type nodeCache struct {
	data       map[string][]NodeData // Clé: hash des endpoints, valeur: données des nœuds
	timestamps map[string]time.Time  // Clé: hash des endpoints, valeur: timestamp de mise en cache
	mutex      sync.RWMutex          // Mutex pour l'accès concurrent
}

// Cache global
var cache = nodeCache{
	data:       make(map[string][]NodeData),
	timestamps: make(map[string]time.Time),
}

// Stockage de la dernière requête
var lastReq = lastRequest{}

// autoRefreshCache rafraîchit périodiquement le cache en arrière-plan
func autoRefreshCache() {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		// Vérifier si nous avons une requête récente à rafraîchir
		lastReq.mutex.RLock()
		names := lastReq.names
		endpoints := lastReq.endpoints
		hasLastRequest := len(names) > 0 && len(endpoints) > 0
		lastReq.mutex.RUnlock()

		if !hasLastRequest {
			log.Println("No recent request to refresh cache")
			continue
		}

		// Forcer le rafraîchissement en ignorant le cache
		refreshedNodes := fetchNodesDataForce(names, endpoints)
		if len(refreshedNodes) > 0 {
			log.Printf("Successfully refreshed cache for %d nodes\n", len(endpoints))
		} else {
			log.Println("Cache refresh failed")
		}
	}
}

// Génère une clé unique pour le cache basée sur les endpoints et les noms
func generateCacheKey(names, endpoints []string) string {
	// Simple concaténation des noms et endpoints
	key := ""
	for i := range names {
		key += names[i] + ":" + endpoints[i] + "|"
	}
	return key
}
