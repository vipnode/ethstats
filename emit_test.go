package main

import (
	"encoding/json"
	"testing"
)

const firstMsg = `{
  "emit": [
    "hello",
    {
      "id": "foo",
      "info": {
        "name": "foo",
        "node": "Geth/v1.8.3-unstable/linux-amd64/go1.10",
        "port": 30303,
        "net": "1",
        "protocol": "les/2",
        "api": "No",
        "os": "linux",
        "os_v": "amd64",
        "client": "0.1.1",
        "canUpdateHistory": true
      },
      "secret": ""
    }
  ]
}`

func TestClientParse(t *testing.T) {
	r, err := parseAuthMsg([]byte(firstMsg))
	if err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	if r.ID != "foo" {
		t.Errorf("incorrect ID: %q", r.ID)
	}
}

func TestEmitMarshal(t *testing.T) {
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
