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

func renderJSON(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

type Server struct {
	collector
	Name stats.ID
}

func (srv *Server) APIHandler(w http.ResponseWriter, r *http.Request) {
	nodeID := r.FormValue("node")

	if nodeID == "" {
		response := struct {
			Nodes []stats.ID `json:"nodes"`
		}{
			Nodes: srv.List(),
		}
		renderJSON(w, response)
		return
	}

	node, ok := srv.Get(stats.ID(nodeID))
	if !ok {
		http.NotFound(w, r)
		return
	}

	renderJSON(w, node)
}

func (srv *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("connected, upgrading", r)
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	go func() {
		err := srv.Join(conn)
		// TODO: Detect normal disconnects and ignore them?
		if err != nil {
			log.Printf("closing connection after error: %s", err)
		}
	}()
}

// Join runs the event loop for a connection, should be run in a goroutine.
func (srv *Server) Join(conn net.Conn) error {
	defer conn.Close()
	var err error
	var emit EmitMessage

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
			report := stats.AuthReport{}
			if err = json.Unmarshal(emit.Payload, &report); err != nil {
				break
			}
			if err = srv.Collect(report); err != nil {
				break
			}
			defer func() {
				srv.Collect(stats.DisconnectReport{report.NodeID()})
			}()
			err = encoder.Encode(&EmitMessage{
				Topic: "ready",
			})
		case "node-ping":
			// Every ethstats implementation ignores the clientTime in
			// the response here, and there is no standard format (eg.
			// geth sends a monotonic offset) so we'll ignore it too.
			sendPayload, err := json.Marshal(&stats.PingReport{srv.Name, time.Now()})
			if err != nil {
				break
			}
			// TODO: We could reuse a sendPayload buffer above
			err = encoder.Encode(&EmitMessage{
				Topic:   "node-pong",
				Payload: sendPayload,
			})
		case "latency":
			report := stats.LatencyReport{}
			if err = json.Unmarshal(emit.Payload, &report); err != nil {
				break
			}
			err = srv.Collect(report)
		case "block":
			report := stats.BlockReport{}
			if err = json.Unmarshal(emit.Payload, &report); err != nil {
				break
			}
			err = srv.Collect(report)
		case "pending":
			report := stats.PendingReport{}
			if err = json.Unmarshal(emit.Payload, &report); err != nil {
				break
			}
			err = srv.Collect(report)
		case "stats":
			report := stats.StatusReport{}
			if err = json.Unmarshal(emit.Payload, &report); err != nil {
				break
			}
			err = srv.Collect(report)
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
