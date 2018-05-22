package main

import "github.com/vipnode/ethstats/stats"

type Node struct {
	Auth    stats.Auth
	Latency stats.Latency
	Block   stats.Block
	Pending stats.Pending
	Status  stats.Status
}
