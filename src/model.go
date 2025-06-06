package main

// NodeData représente les données d'un nœud Storj
type NodeData struct {
	Name          string           `json:"name"`
	BandwidthData []BandwidthDaily `json:"bandwidthData"`
	StorageData   []StorageDaily   `json:"storageData"`
	Earnings      EarningsData     `json:"earnings"`
}

// BandwidthDaily représente les données d'utilisation quotidienne de bande passante
type BandwidthDaily struct {
	IntervalStart string         `json:"intervalStart"`
	Egress        BandwidthUsage `json:"egress"`
	Ingress       BandwidthUsage `json:"ingress"`
}

// BandwidthUsage détaille l'utilisation de bande passante
type BandwidthUsage struct {
	Usage  float64 `json:"usage"`
	Repair float64 `json:"repair"`
	Audit  float64 `json:"audit"`
}

// StorageDaily représente les données quotidiennes de stockage
type StorageDaily struct {
	IntervalStart    string  `json:"intervalStart"`
	AtRestTotalBytes float64 `json:"atRestTotalBytes"`
}

// DashboardData contient toutes les données pour le tableau de bord
type DashboardData struct {
	Nodes           []NodeData      `json:"nodes"`
	EarningsPayouts EarningsPayouts `json:"earningsPayouts"`
	Error           string          `json:"error,omitempty"`
}

// EarningsData représente les données de revenus d'un nœud Storj
type EarningsData struct {
	CurrentMonth struct {
		EgressBandwidth         float64 `json:"egressBandwidth"`
		EgressBandwidthPayout   float64 `json:"egressBandwidthPayout"`
		EgressRepairAudit       float64 `json:"egressRepairAudit"`
		EgressRepairAuditPayout float64 `json:"egressRepairAuditPayout"`
		DiskSpace               float64 `json:"diskSpace"`
		DiskSpacePayout         float64 `json:"diskSpacePayout"`
		HeldRate                float64 `json:"heldRate"`
		Payout                  float64 `json:"payout"`
		Held                    float64 `json:"held"`
	} `json:"currentMonth"`
	PreviousMonth struct {
		EgressBandwidth         float64 `json:"egressBandwidth"`
		EgressBandwidthPayout   float64 `json:"egressBandwidthPayout"`
		EgressRepairAudit       float64 `json:"egressRepairAudit"`
		EgressRepairAuditPayout float64 `json:"egressRepairAuditPayout"`
		DiskSpace               float64 `json:"diskSpace"`
		DiskSpacePayout         float64 `json:"diskSpacePayout"`
		HeldRate                float64 `json:"heldRate"`
		Payout                  float64 `json:"payout"`
		Held                    float64 `json:"held"`
	} `json:"previousMonth"`
	CurrentMonthExpectations float64 `json:"currentMonthExpectations"`
}

type EarningsPayouts struct {
	EgressBandwidthPayout    float64 `json:"egressBandwidthPayout"`
	EgressRepairAuditPayout  float64 `json:"egressRepairAuditPayout"`
	DiskSpacePayout          float64 `json:"diskSpacePayout"`
	CurrentMonthExpectations float64 `json:"currentMonthExpectations"`
	CurrentMonthTotal        float64 `json:"currentMonthTotal"`
	PreviousMonthTotal       float64 `json:"previousMonthTotal"`
	TotalHeld                float64 `json:"totalHeld"`
}
