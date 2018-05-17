package main

import (
	"encoding/json"
	"fmt"

	"github.com/vipnode/ethstats/stats"
)

func parseAuthMsg(data []byte) (*stats.AuthMsg, error) {
	var emitMsg EmitMessage
	if err := emitMsg.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	if emitMsg.Topic != "hello" {
		return nil, fmt.Errorf("unexpected emit topic: %q", emitMsg.Topic)
	}

	var authMsg stats.AuthMsg
	if err := json.Unmarshal(emitMsg.Payload, &authMsg); err != nil {
		return nil, err
	}
	return &authMsg, nil
}
