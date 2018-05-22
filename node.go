package main

import "github.com/vipnode/ethstats/stats"

type Node struct {
	Auth         stats.Auth
	Latency      stats.LatencyReport
	BlockStats   stats.Block
	PendingStats stats.Pending
	NodeStats    stats.Status
}
