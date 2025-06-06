package main

import (
	"fmt"
	"net/http"
	"time"
)

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Extraire les paramètres de la requête au nouveau format
	// Par exemple: ?Node1=192.168.1.30:14002&Node2=192.168.1.30:14003
	var names []string
	var endpoints []string

	// Parcourir tous les paramètres de l'URL
	for name, values := range r.URL.Query() {
		if len(values) > 0 && values[0] != "" {
			names = append(names, name)
			endpoints = append(endpoints, values[0])
		}
	}

	// Vérifier que nous avons des endpoints valides
	if len(names) == 0 || len(endpoints) == 0 {
		DashboardPage(DashboardData{
			Error: "Please provide at least one node endpoint in the query, ex : ?NodeName=endpoint:port",
		}).Render(r.Context(), w)
		return
	}

	// Sauvegarder cette requête pour le rafraîchissement automatique
	lastReq.mutex.Lock()
	lastReq.names = make([]string, len(names))
	lastReq.endpoints = make([]string, len(endpoints))
	copy(lastReq.names, names)
	copy(lastReq.endpoints, endpoints)
	lastReq.timestamp = time.Now()
	lastReq.mutex.Unlock()

	// Récupérer les données des nœuds (avec cache)
	nodes := fetchNodesData(names, endpoints)

	// Vérifier si nous avons des données
	if len(nodes) == 0 {
		DashboardPage(DashboardData{
			Error: "No nodes found or data could not be fetched. Please check your endpoints.",
		}).Render(r.Context(), w)
		return
	}

	// Initialiser les variables pour les totaux
	var earningPayout = EarningsPayouts{
		EgressBandwidthPayout:    0,
		EgressRepairAuditPayout:  0,
		DiskSpacePayout:          0,
		CurrentMonthExpectations: 0,
		CurrentMonthTotal:        0,
		PreviousMonthTotal:       0,
		TotalHeld:                0,
	}

	// Calculer les totaux à partir de tous les nœuds
	for _, node := range nodes {
		// Current month: somme des payouts par catégorie
		earningPayout.CurrentMonthTotal += node.Earnings.CurrentMonth.EgressBandwidthPayout +
			node.Earnings.CurrentMonth.EgressRepairAuditPayout +
			node.Earnings.CurrentMonth.DiskSpacePayout

		// Previous month: somme des payouts par catégorie
		earningPayout.PreviousMonthTotal += node.Earnings.PreviousMonth.EgressBandwidthPayout +
			node.Earnings.PreviousMonth.EgressRepairAuditPayout +
			node.Earnings.PreviousMonth.DiskSpacePayout

		// Estimation pour le mois
		earningPayout.CurrentMonthExpectations += node.Earnings.CurrentMonthExpectations

		// Ajouter les payouts individuels
		earningPayout.EgressBandwidthPayout += node.Earnings.CurrentMonth.EgressBandwidthPayout +
			node.Earnings.PreviousMonth.EgressBandwidthPayout
		earningPayout.EgressRepairAuditPayout += node.Earnings.CurrentMonth.EgressRepairAuditPayout +
			node.Earnings.PreviousMonth.EgressRepairAuditPayout
		earningPayout.DiskSpacePayout += node.Earnings.CurrentMonth.DiskSpacePayout +
			node.Earnings.PreviousMonth.DiskSpacePayout

		// Total held
		earningPayout.TotalHeld += node.Earnings.CurrentMonth.Held + node.Earnings.PreviousMonth.Held
	}

	// Rendre le template avec les données
	DashboardPage(DashboardData{
		Nodes:           nodes,
		EarningsPayouts: earningPayout,
	}).Render(r.Context(), w)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok","service":"storj-multinode-dashboard"}`)
}
