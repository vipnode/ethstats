package main

import (
	"encoding/json"
	"testing"
)

func TestEmitMessage(t *testing.T) {
	tests := []struct {
		Emit     EmitMessage
		Expected string
	}{
		{
			EmitMessage{
				Topic:   "hello",
				Payload: json.RawMessage(`{"foo": 42}`),
			},
			`{"emit":["hello",{"foo":42}]}`,
		},
		{
			EmitMessage{
				Topic: "ack",
			},
			`{"emit":["ack"]}`,
		},
	}

	for _, tc := range tests {
		out, err := tc.Emit.MarshalJSON()
		if err != nil {
			t.Errorf("failed to marshal: %q", err)
		}

		if got, want := string(out), tc.Expected; got != want {
			t.Errorf("got:\n\t%s; want\n\t%s", got, want)
		}
	}
}

func TestMarshalEmit(t *testing.T) {
	out, err := MarshalEmit("ready", nil)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(out), `{"emit":["ready"]}`; got != want {
		t.Errorf("got:\n\t%s; want\n\t%s", got, want)
	}
}
