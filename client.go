package main

import (
	"encoding/json"

	"github.com/vipnode/ethstats/stats"
)

type Client struct {
}

func parseAuthMsg(msg []byte) (string, *stats.AuthMsg, error) {
	var authMsg *stats.AuthMsg
	if err := json.Unmarshal(msg, authMsg); err != nil {
		return "", nil, err
	}
	return "foo", authMsg, nil
}
