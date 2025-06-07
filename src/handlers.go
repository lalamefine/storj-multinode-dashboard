package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var names []string
	var endpoints []string

	urlParams := r.URL.Query()
	nodesEnv := os.Getenv("NODES")

	if nodesEnv != "" {
		nodePairs := strings.Split(nodesEnv, ",")
		for _, pair := range nodePairs {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
				names = append(names, parts[0])
				endpoints = append(endpoints, parts[1])
			}
		}
	} else {
		for name, values := range urlParams {
			if len(values) > 0 && values[0] != "" {
				names = append(names, name)
				endpoints = append(endpoints, values[0])
			}
		}
	}

	if len(names) == 0 || len(endpoints) == 0 {
		DashboardPage(DashboardData{
			Error: "Please provide at least one node endpoint either in the URL query or via the NODES environment variable.",
		}).Render(r.Context(), w)
		return
	}

	lastReq.mutex.Lock()
	lastReq.names = make([]string, len(names))
	lastReq.endpoints = make([]string, len(endpoints))
	copy(lastReq.names, names)
	copy(lastReq.endpoints, endpoints)
	lastReq.timestamp = time.Now()
	lastReq.mutex.Unlock()

	nodes := fetchNodesData(names, endpoints)

	if len(nodes) == 0 {
		DashboardPage(DashboardData{
			Error: "No nodes found or data could not be fetched. Please check your endpoints.",
		}).Render(r.Context(), w)
		return
	}

	var earningPayout = EarningsPayouts{
		EgressBandwidthPayout:    0,
		EgressRepairAuditPayout:  0,
		DiskSpacePayout:          0,
		CurrentMonthExpectations: 0,
		CurrentMonthTotal:        0,
		PreviousMonthTotal:       0,
		TotalHeld:                0,
	}

	for _, node := range nodes {
		earningPayout.CurrentMonthTotal += node.Earnings.CurrentMonth.EgressBandwidthPayout +
			node.Earnings.CurrentMonth.EgressRepairAuditPayout +
			node.Earnings.CurrentMonth.DiskSpacePayout

		earningPayout.PreviousMonthTotal += node.Earnings.PreviousMonth.EgressBandwidthPayout +
			node.Earnings.PreviousMonth.EgressRepairAuditPayout +
			node.Earnings.PreviousMonth.DiskSpacePayout

		earningPayout.CurrentMonthExpectations += node.Earnings.CurrentMonthExpectations

		earningPayout.EgressBandwidthPayout += node.Earnings.CurrentMonth.EgressBandwidthPayout +
			node.Earnings.PreviousMonth.EgressBandwidthPayout
		earningPayout.EgressRepairAuditPayout += node.Earnings.CurrentMonth.EgressRepairAuditPayout +
			node.Earnings.PreviousMonth.EgressRepairAuditPayout
		earningPayout.DiskSpacePayout += node.Earnings.CurrentMonth.DiskSpacePayout +
			node.Earnings.PreviousMonth.DiskSpacePayout

		earningPayout.TotalHeld += node.Earnings.CurrentMonth.Held + node.Earnings.PreviousMonth.Held
	}

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
