package main

import (
	"encoding/json"
	"fmt"
)

// EmitMessage contains a parsed SocksJS-style pubsub event emit.
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
	if len(msg.Emit) == 0 {
		return fmt.Errorf("missing emit fields")
	}
	if err := json.Unmarshal(msg.Emit[0], &emit.Topic); err != nil {
		return err
	}
	if len(msg.Emit) > 1 {
		emit.Payload = msg.Emit[1]
	}
	return nil
}

func (emit *EmitMessage) MarshalJSON() ([]byte, error) {
	msg := struct {
		Emit []json.RawMessage `json:"emit"`
	}{}

	if emit.Topic == "" {
		return nil, fmt.Errorf("missing topic")
	}
	rawTopic, err := json.Marshal(emit.Topic)
	if err != nil {
		return nil, err
	}

	msg.Emit = append(msg.Emit, rawTopic)
	if emit.Payload != nil {
		msg.Emit = append(msg.Emit, emit.Payload)
	}
	return json.Marshal(msg)
}

func MarshalEmit(topic string, payload interface{}) ([]byte, error) {
	emit := EmitMessage{
		Topic: topic,
	}

	if payload != nil {
		rawPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		emit.Payload = rawPayload
	}

	return emit.MarshalJSON()
}
