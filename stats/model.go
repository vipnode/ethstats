package stats // import "github.com/vipnode/ethstats/stats"

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// These structs are mostly borrowed from github.com/ethereu/go-ethereum/ethstats

// NodeInfo is the collection of metainformation about a node that is displayed
// on the monitoring page.
type NodeInfo struct {
	Name     string `json:"name"`
	Node     string `json:"node"`
	Port     int    `json:"port"`
	Network  string `json:"net"`
	Protocol string `json:"protocol"`
	API      string `json:"api"`
	Os       string `json:"os"`
	OsVer    string `json:"os_v"`
	Client   string `json:"client"`
	History  bool   `json:"canUpdateHistory"`
}

// AuthMsg is the authentication infos needed to login to a monitoring server.
// TODO: Rename to Auth?
type AuthMsg struct {
	ID     string   `json:"id"`
	Info   NodeInfo `json:"info"`
	Secret string   `json:"secret"`
}

// BlockStats is the information to report about individual blocks.
// TODO: Rename to stats.Block?
type BlockStats struct {
	Number     *big.Int       `json:"number"`
	Hash       common.Hash    `json:"hash"`
	ParentHash common.Hash    `json:"parentHash"`
	Timestamp  *big.Int       `json:"timestamp"`
	Miner      common.Address `json:"miner"`
	GasUsed    uint64         `json:"gasUsed"`
	GasLimit   uint64         `json:"gasLimit"`
	Diff       string         `json:"difficulty"`
	TotalDiff  string         `json:"totalDifficulty"`
	Txs        []TxStats      `json:"transactions"`
	TxHash     common.Hash    `json:"transactionsRoot"`
	Root       common.Hash    `json:"stateRoot"`
	Uncles     uncleStats     `json:"uncles"`
}

// uncleStats is a custom wrapper around an uncle array to force serializing
// empty arrays instead of returning null for them.
type uncleStats []*types.Header

func (s uncleStats) MarshalJSON() ([]byte, error) {
	if uncles := ([]*types.Header)(s); len(uncles) > 0 {
		return json.Marshal(uncles)
	}
	return []byte("[]"), nil
}

// TxStats is the information to report about individual transactions.
type TxStats struct {
	Hash common.Hash `json:"hash"`
}

// PendingStats is the information to report about pending transactions.
// TODO: Rename to stats.Pending?
type PendingStats struct {
	Pending int `json:"pending"`
}

// NodeStats is the information to report about the local node.
// TODO: Rename to stats.Node?
type NodeStats struct {
	Active   bool `json:"active"`
	Syncing  bool `json:"syncing"`
	Mining   bool `json:"mining"`
	Hashrate int  `json:"hashrate"`
	Peers    int  `json:"peers"`
	GasPrice int  `json:"gasPrice"`
	Uptime   int  `json:"uptime"`
}

type StatsReport struct {
	ID    string    `json:"id"`
	Stats NodeStats `json:"stats"`
}

type NodePing struct {
	ID         string    `json:"id"`
	ClientTime time.Time `json:"clientTime"`
}

type NodeLatency struct {
	ID      string `json:"id"`
	Latency string `json:"latency"`
}
