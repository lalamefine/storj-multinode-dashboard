package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func fetchNodesDataForce(names, endpoints []string) []NodeData {
	nodes := fetchNodesDataInternal(names, endpoints)

	if len(nodes) > 0 {
		cacheKey := generateCacheKey(names, endpoints)
		cache.mutex.Lock()
		cache.data[cacheKey] = nodes
		cache.timestamps[cacheKey] = time.Now()
		cache.mutex.Unlock()
	}

	return nodes
}

func fetchNodesData(names, endpoints []string) []NodeData {
	cacheKey := generateCacheKey(names, endpoints)

	cache.mutex.RLock()
	cachedData, exists := cache.data[cacheKey]
	cachedTime, timeExists := cache.timestamps[cacheKey]
	cache.mutex.RUnlock()

	if exists && timeExists && time.Since(cachedTime) < cacheExpiration {
		log.Printf("Data fetched from cache for %d nodes", len(endpoints))
		return cachedData
	}

	return fetchNodesDataForce(names, endpoints)
}

func fetchNodesDataInternal(names, endpoints []string) []NodeData {
	log.Printf("Fetching fresh data for %d nodes", len(endpoints))
	nodes := make([]NodeData, len(endpoints))
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var wg sync.WaitGroup
	for i := range endpoints {
		wg.Add(1)

		go func(index int, nodeName, nodeEndpoint string) {
			defer wg.Done()
			node, success := fetchSingleNodeData(client, nodeName, nodeEndpoint)
			if success {
				nodes[index] = node
			}
		}(i, names[i], endpoints[i])
	}

	wg.Wait()
	var validNodes []NodeData
	for _, node := range nodes {
		if node.Name != "" {
			validNodes = append(validNodes, node)
		}
	}

	return validNodes
}

func fetchSingleNodeData(client *http.Client, name, endpoint string) (NodeData, bool) {
	baseUrl := "http://" + endpoint
	url := baseUrl + "/api/sno/satellites"
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error connecting to %s: %v", endpoint, err)
		return NodeData{}, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Invalid response status for %s: %d", endpoint, resp.StatusCode)
		return NodeData{}, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from %s: %v", endpoint, err)
		return NodeData{}, false
	}

	var data struct {
		BandwidthDaily []BandwidthDaily `json:"bandwidthDaily"`
		StorageDaily   []StorageDaily   `json:"storageDaily"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error decoding JSON response from %s: %v", endpoint, err)
		return NodeData{}, false
	}

	var earnings EarningsData
	earningsUrl := baseUrl + "/api/sno/estimated-payout"
	respEarnings, err := client.Get(earningsUrl)

	if err != nil {
		log.Printf("Error connecting to earnings endpoint %s: %v", earningsUrl, err)
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

	return NodeData{
		Name:          name,
		BandwidthData: data.BandwidthDaily,
		StorageData:   data.StorageDaily,
		Earnings:      earnings,
	}, true
}
