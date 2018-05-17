package main

import "testing"

const authMsg = `{
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

const pingMsg = `{
  "emit": [
    "node-ping",
    {
      "clientTime": "2018-05-17T21:15:15.389Z",
      "id": "foo"
    }
  ]
}`

// NOTE: For some reason geth uses time.Time.string() which includes a monotonic offset
// that does not unmarshal properly.
// Eg: {"clientTime": "2018-05-17 16:53:43.96985387 -0400 EDT m=+15.136170456"}

func TestClientParseAuth(t *testing.T) {
	var emitMsg EmitMessage
	if err := emitMsg.UnmarshalJSON([]byte(authMsg)); err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	r, err := parseAuthMsg(emitMsg)
	if err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	if r.ID != "foo" {
		t.Errorf("incorrect ID: %q", r.ID)
	}
}

func TestClientParsePing(t *testing.T) {
	var emitMsg EmitMessage
	if err := emitMsg.UnmarshalJSON([]byte(pingMsg)); err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	r, err := parsePingMsg(emitMsg)
	if err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	if r.ID != "foo" {
		t.Errorf("incorrect ID: %q", r.ID)
	}

	if r.ClientTime.Second() != 15 {
		t.Errorf("incorrect timestamp: %q", r.ClientTime)
	}
}
