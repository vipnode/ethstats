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

	for {
		// TODO: Reuse byte buffer
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			return err
		}
		log.Print("received", string(msg), op, err)

		emit := EmitMessage{}
		if err := json.Unmarshal(msg, &emit); err != nil {
			return err
		}

		var out []byte
		switch topic := emit.Topic; topic {
		case "hello":
			out, err = MarshalEmit("ready", nil)
		case "node-ping":
			// Every ethstats implementation ignores the clientTime in
			// the response here, and there is no standard format (eg.
			// geth sends a monotonic offset) so we'll ignore it too.
			out, err = MarshalEmit(
				"node-pong",
				stats.NodePing{srv.Name, time.Now()},
			)
		default:
			continue
		}

		if err != nil {
			return err
		}
		if err := wsutil.WriteServerMessage(conn, op, out); err != nil {
			return err
		}

		log.Print("sent", string(out), op, err)
	}
}
