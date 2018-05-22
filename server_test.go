package main

import (
	"encoding/json"
	"testing"

	"github.com/vipnode/ethstats/stats"
)

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

func TestParseAuth(t *testing.T) {
	var emitMsg EmitMessage
	if err := emitMsg.UnmarshalJSON([]byte(authMsg)); err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	if emitMsg.Topic != "hello" {
		t.Errorf("unexpected emit topic: %q", emitMsg.Topic)
	}

	node := Node{}
	if err := json.Unmarshal(emitMsg.Payload, &node.Auth); err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	if node.Auth.ID != "foo" {
		t.Errorf("incorrect ID: %q", node.Auth.ID)
	}
}

func TestParsePing(t *testing.T) {
	var emitMsg EmitMessage
	if err := emitMsg.UnmarshalJSON([]byte(pingMsg)); err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	if emitMsg.Topic != "node-ping" {
		t.Errorf("unexpected emit topic: %q", emitMsg.Topic)
	}

	var r stats.NodePing
	if err := json.Unmarshal(emitMsg.Payload, &r); err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	if r.ID != "foo" {
		t.Errorf("incorrect ID: %q", r.ID)
	}

	if r.ClientTime.Second() != 15 {
		t.Errorf("incorrect timestamp: %q", r.ClientTime)
	}
}

const blockMsg = `{"emit":["block",{"block":{"number":5273251,"hash":"0xe11bc629a85375753ba5a043e5b44c05dedbdb484ed8956f9aec07bf3d93fde5","parentHash":"0x10aa19d73522d15cf004ca602b3b87d79bb903d5f7ba8745fc7959534047c7de","timestamp":1521317517,"miner":"0xb2930b35844a230f00e51431acae96fe543a0347","gasUsed":7984834,"gasLimit":7999992,"difficulty":"3291915733734816","totalDifficulty":"3102951517281028058241","transactions":[],"transactionsRoot":"0xe2fdfcc5707a06727f7624ae01c8a7128194b4fc88579375f2ab96e3bdc12d08","stateRoot":"0xdb34c6952061b45c9f4875ed70475cddd3cee0ba016afbd4c2418bbe9ca539d4","uncles":[]},"id":"a"}]}`

func TestParseBlock(t *testing.T) {
	var emitMsg EmitMessage
	if err := emitMsg.UnmarshalJSON([]byte(blockMsg)); err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	if emitMsg.Topic != "block" {
		t.Errorf("unexpected emit topic: %q", emitMsg.Topic)
	}

	node := Node{}
	container := struct {
		Block *stats.BlockStats `json:"block"`
		ID    string            `json:"id"`
	}{
		Block: &node.BlockStats,
	}
	if err := json.Unmarshal(emitMsg.Payload, &container); err != nil {
		t.Fatalf("failed to parse: %q", err)
	}

	if node.BlockStats.Number.String() != "5273251" {
		t.Errorf("incorrect block number: %q", node.BlockStats.Number)
	}

	if node.BlockStats.Hash.String() != "0xe11bc629a85375753ba5a043e5b44c05dedbdb484ed8956f9aec07bf3d93fde5" {
		t.Errorf("incorrect block hash: %q", node.BlockStats.Hash)
	}
}

const statsMsg = `{"emit":["stats",{"id":"a","stats":{"active":true,"syncing":true,"mining":false,"hashrate":0,"peers":0,"gasPrice":0,"uptime":100}}]}`
