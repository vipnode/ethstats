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
	node := Node{}

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			return err
		}
		log.Print("received", string(msg), op, err)

		emit := EmitMessage{}
		if err := json.Unmarshal(msg, &emit); err != nil {
			return err
		}

		// TODO: Support relaying by trusting and mapping ID?
		// TODO: Reuse out buffer
		var out []byte
		switch topic := emit.Topic; topic {
		case "hello":
			if err = json.Unmarshal(emit.Payload, &node.Auth); err != nil {
				break
			}
			out, err = MarshalEmit("ready", nil)
		case "node-ping":
			// Every ethstats implementation ignores the clientTime in
			// the response here, and there is no standard format (eg.
			// geth sends a monotonic offset) so we'll ignore it too.
			out, err = MarshalEmit(
				"node-pong",
				stats.NodePing{srv.Name, time.Now()},
			)
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
			log.Printf("error %q on message: %s", err, msg)
			return err
		}
		if len(out) == 0 {
			continue
		}
		if err := wsutil.WriteServerMessage(conn, op, out); err != nil {
			return err
		}

		log.Print("sent", string(out), op, err)
	}
}
