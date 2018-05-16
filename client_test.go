package main

import "testing"

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

func TestClient(t *testing.T) {
	cmd, r, err := parseAuthMsg([]byte(firstMsg))
	if err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	if cmd != "hello" {
		t.Errorf("incorrect command: %q", cmd)
	}

	if r.ID != "foo" {
		t.Errorf("incorrect ID: %q", r.ID)
	}
}
