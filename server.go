package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vipnode/ethstats/stats"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // We don't care about XSS
}

type Server struct {
	Name string
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("connected, upgrading", r)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
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
func (srv *Server) Join(conn *websocket.Conn) error {
	defer conn.Close()
	var err error

	node := Node{}
	emit := EmitMessage{}

	for {
		// Decode next message
		if err = conn.ReadJSON(&emit); err != nil {
			return err
		}

		log.Printf("received from %s: %s", conn.RemoteAddr(), emit.Topic)

		// TODO: Support relaying by trusting and mapping ID?
		switch topic := emit.Topic; topic {
		case "hello":
			if err = json.Unmarshal(emit.Payload, &node.Auth); err != nil {
				break
			}
			// TODO: Assert ID?
			err = conn.WriteJSON(&EmitMessage{
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
			err = conn.WriteJSON(&EmitMessage{
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
	}
}
