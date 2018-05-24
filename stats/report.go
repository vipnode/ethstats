package stats

import "time"

// Report is a container for some stats about a node.
type Report interface {
	NodeID() ID
}

// AuthReport contains the authorization needed to log into a monitoring server.
type AuthReport struct {
	ID     ID     `json:"id"`
	Info   Info   `json:"info"`
	Secret string `json:"secret"`
}

func (r AuthReport) NodeID() ID {
	return r.ID
}

// StatusReport contains the Status for a specific node ID
type StatusReport struct {
	ID     ID     `json:"id"`
	Status Status `json:"stats"`
}

func (r StatusReport) NodeID() ID {
	return r.ID
}

// PendingReport contains the Pending stats for a specific node ID
type PendingReport struct {
	ID      ID      `json:"id"`
	Pending Pending `json:"stats"`
}

func (r PendingReport) NodeID() ID {
	return r.ID
}

// BlockReport contains the Block stats for a specific node ID
type BlockReport struct {
	ID    ID    `json:"id"`
	Block Block `json:"block"`
}

func (r BlockReport) NodeID() ID {
	return r.ID
}

// PingReport contains the client time for a specific node ID
type PingReport struct {
	ID         ID        `json:"id"`
	ClientTime time.Time `json:"clientTime"`
}

func (r PingReport) NodeID() ID {
	return r.ID
}

// LatencyReport contains the latency to a specific node ID
type LatencyReport struct {
	ID      ID      `json:"id"`
	Latency Latency `json:"latency"`
}

func (r LatencyReport) NodeID() ID {
	return r.ID
}

// DisconnectReport signals a disconnect event for a specific node ID.
type DisconnectReport struct {
	ID ID `json:"id"`
}

func (r DisconnectReport) NodeID() ID {
	return r.ID
}
