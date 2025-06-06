package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// fetchNodesDataForce force le rafraîchissement du cache sans vérifier s'il est valide
func fetchNodesDataForce(names, endpoints []string) []NodeData {
	nodes := fetchNodesDataInternal(names, endpoints)

	// Mettre à jour le cache si nous avons récupéré des données
	if len(nodes) > 0 {
		cacheKey := generateCacheKey(names, endpoints)
		cache.mutex.Lock()
		cache.data[cacheKey] = nodes
		cache.timestamps[cacheKey] = time.Now()
		cache.mutex.Unlock()
	}

	return nodes
}

// fetchNodesData récupère les données des nœuds Storj (avec mise en cache)
func fetchNodesData(names, endpoints []string) []NodeData {
	cacheKey := generateCacheKey(names, endpoints)

	// Vérifier si les données sont en cache et toujours valides
	cache.mutex.RLock()
	cachedData, exists := cache.data[cacheKey]
	cachedTime, timeExists := cache.timestamps[cacheKey]
	cache.mutex.RUnlock()

	// Si les données sont en cache et encore valides
	if exists && timeExists && time.Since(cachedTime) < cacheExpiration {
		log.Printf("Data fetched from cache for %d nodes", len(endpoints))
		return cachedData
	}

	return fetchNodesDataForce(names, endpoints)
}

// fetchNodesDataInternal contient la logique de récupération des données
func fetchNodesDataInternal(names, endpoints []string) []NodeData {
	// Si les données ne sont pas en cache ou sont expirées, les récupérer
	log.Printf("Fetching fresh data for %d nodes", len(endpoints))
	var nodes []NodeData

	// Créer un client HTTP avec timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Récupérer les données pour chaque nœud
	for i, endpoint := range endpoints {
		name := names[i]
		baseUrl := "http://" + endpoint

		// Récupérer les données de bande passante et stockage
		url := baseUrl + "/api/sno/satellites"

		// Faire la requête
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("Error connecting to %s: %v", endpoint, err)
			continue
		}
		defer resp.Body.Close()

		// Vérifier le statut de la réponse
		if resp.StatusCode != http.StatusOK {
			log.Printf("Invalid response status for %s: %d", endpoint, resp.StatusCode)
			continue
		}

		// Lire le corps de la réponse
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body from %s: %v", endpoint, err)
			continue
		}

		// Décoder la réponse JSON
		var data struct {
			BandwidthDaily []BandwidthDaily `json:"bandwidthDaily"`
			StorageDaily   []StorageDaily   `json:"storageDaily"`
		}

		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("Error decoding JSON response from %s: %v", endpoint, err)
			continue
		}

		// Récupérer les données de revenus
		earningsUrl := baseUrl + "/api/sno/estimated-payout"
		respEarnings, err := client.Get(earningsUrl)
		var earnings EarningsData
		if err != nil {
			log.Printf("Error connecting to earnings endpoint %s: %v", earningsUrl, err)
			// Continuer avec les autres données
		} else {
			defer respEarnings.Body.Close()

			if respEarnings.StatusCode == http.StatusOK {
				bodyEarnings, err := io.ReadAll(respEarnings.Body)
				if err == nil {
					if err := json.Unmarshal(bodyEarnings, &earnings); err != nil {
						log.Printf("Error decoding earnings JSON from %s: %v", endpoint, err)
					}
				}
			}
		}

		// Ajouter les données du nœud
		nodes = append(nodes, NodeData{
			Name:          name,
			BandwidthData: data.BandwidthDaily,
			StorageData:   data.StorageDaily,
			Earnings:      earnings,
		})
	}

	return nodes
}
