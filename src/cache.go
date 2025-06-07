package main

import (
	"log"
	"sync"
	"time"
)

type lastRequest struct {
	names     []string
	endpoints []string
	timestamp time.Time
	mutex     sync.RWMutex
}

type nodeCache struct {
	data       map[string][]NodeData
	timestamps map[string]time.Time
	mutex      sync.RWMutex
}

var cache = nodeCache{
	data:       make(map[string][]NodeData),
	timestamps: make(map[string]time.Time),
}

var lastReq = lastRequest{}

func autoRefreshCache() {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		lastReq.mutex.RLock()
		names := lastReq.names
		endpoints := lastReq.endpoints
		hasLastRequest := len(names) > 0 && len(endpoints) > 0
		lastReq.mutex.RUnlock()

		if !hasLastRequest {
			log.Println("No recent request to refresh cache")
			continue
		}

		refreshedNodes := fetchNodesDataForce(names, endpoints)
		if len(refreshedNodes) > 0 {
			log.Printf("Successfully refreshed cache for %d nodes\n", len(endpoints))
		} else {
			log.Println("Cache refresh failed")
		}
	}
}

func generateCacheKey(names, endpoints []string) string {
	key := ""
	for i := range names {
		key += names[i] + ":" + endpoints[i] + "|"
	}
	return key
}
