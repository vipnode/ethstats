package main

import "github.com/vipnode/ethstats/stats"

type Node struct {
	Auth         stats.AuthMsg
	Latency      stats.NodeLatency
	BlockStats   stats.BlockStats
	PendingStats stats.PendingStats
	NodeStats    stats.NodeStats
}
