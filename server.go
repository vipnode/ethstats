package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/vipnode/ethstats/stats"
)

type Server struct {
	Name string
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("connected, upgrading", r)
	conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	go func() {
		err := srv.Join(conn)
		if err != nil {
			log.Printf("closing connection after error: %s", err)
		}
	}()
}

// Join runs the event loop for a connection, should be run in a goroutine.
func (srv *Server) Join(conn net.Conn) error {
	defer conn.Close()
	var err error

	node := Node{}
	emit := EmitMessage{}

	r := wsutil.NewReader(conn, ws.StateServerSide)
	w := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)

	decoder := json.NewDecoder(r)
	encoder := json.NewEncoder(w)

	for {
		// Prepare for the next message
		if _, err = r.NextFrame(); err != nil {
			return err
		}
		// Decode next message
		if err = decoder.Decode(&emit); err != nil {
			return err
		}

		log.Printf("%s: received topic: %s", conn.RemoteAddr(), emit.Topic)

		// TODO: Support relaying by trusting and mapping ID?
		switch topic := emit.Topic; topic {
		case "hello":
			if err = json.Unmarshal(emit.Payload, &node.Auth); err != nil {
				break
			}
			// TODO: Assert ID?
			err = encoder.Encode(&EmitMessage{
				Topic: "ready",
			})
		case "node-ping":
			// Every ethstats implementation ignores the clientTime in
			// the response here, and there is no standard format (eg.
			// geth sends a monotonic offset) so we'll ignore it too.
			sendPayload, err := json.Marshal(&stats.NodePing{srv.Name, time.Now()})
			if err != nil {
				break
			}
			// TODO: We could reuse a sendPayload buffer above
			err = encoder.Encode(&EmitMessage{
				Topic:   "node-pong",
				Payload: sendPayload,
			})
		case "latency":
			err = json.Unmarshal(emit.Payload, &node.Latency)
		case "block":
			// Contained in {"block": ..., "id": ...}
			container := struct {
				Block *stats.BlockStats `json:"block"`
				ID    string            `json:"id"`
			}{
				Block: &node.BlockStats,
			}
			err = json.Unmarshal(emit.Payload, &container)
		case "pending":
			// Contained in {"stats": ..., "id": ...}
			container := struct {
				Stats *stats.PendingStats `json:"stats"`
				ID    string              `json:"id"`
			}{
				Stats: &node.PendingStats,
			}
			err = json.Unmarshal(emit.Payload, &container)
		case "stats":
			// Contained in {"stats": ..., "id": ...}
			container := struct {
				Stats *stats.NodeStats `json:"stats"`
				ID    string           `json:"id"`
			}{
				Stats: &node.NodeStats,
			}
			err = json.Unmarshal(emit.Payload, &container)
		default:
			continue
		}

		if err != nil {
			return err
		}

		// Write buffer
		if err = w.Flush(); err != nil {
			return err
		}
	}
}
