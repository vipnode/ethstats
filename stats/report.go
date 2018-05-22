package stats

import "time"

type StatusReport struct {
	ID     string `json:"id"`
	Status Status `json:"stats"`
}

type PendingReport struct {
	ID      string  `json:"id"`
	Pending Pending `json:"stats"`
}

type BlockReport struct {
	ID    string `json:"id"`
	Block Block  `json:"block"`
}

type PingReport struct {
	ID         string    `json:"id"`
	ClientTime time.Time `json:"clientTime"`
}

type LatencyReport struct {
	ID      string  `json:"id"`
	Latency Latency `json:"latency"`
}
