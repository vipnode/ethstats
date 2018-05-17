package main

import (
	"encoding/json"
	"fmt"

	"github.com/vipnode/ethstats/stats"
)

type EmitMessage struct {
	Topic   string
	Payload json.RawMessage
}

func (emit *EmitMessage) UnmarshalJSON(data []byte) error {
	msg := struct {
		Emit []json.RawMessage `json:"emit"`
	}{}

	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	if len(msg.Emit) != 2 {
		return fmt.Errorf("expected emit tuple of 2, got: %d", len(msg.Emit))
	}
	if err := json.Unmarshal(msg.Emit[0], &emit.Topic); err != nil {
		return err
	}
	emit.Payload = msg.Emit[1]
	return nil
}

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
