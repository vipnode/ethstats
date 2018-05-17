package main

import (
	"encoding/json"
	"fmt"

	"github.com/vipnode/ethstats/stats"
)

func parseAuthMsg(emitMsg EmitMessage) (*stats.AuthMsg, error) {
	if emitMsg.Topic != "hello" {
		return nil, fmt.Errorf("unexpected emit topic: %q", emitMsg.Topic)
	}

	var authMsg stats.AuthMsg
	if err := json.Unmarshal(emitMsg.Payload, &authMsg); err != nil {
		return nil, err
	}
	return &authMsg, nil
}

func parsePingMsg(emitMsg EmitMessage) (*stats.NodePing, error) {
	if emitMsg.Topic != "node-ping" {
		return nil, fmt.Errorf("unexpected emit topic: %q", emitMsg.Topic)
	}

	var pingMsg stats.NodePing
	if err := json.Unmarshal(emitMsg.Payload, &pingMsg); err != nil {
		return nil, err
	}
	return &pingMsg, nil
}
