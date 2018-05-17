package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/vipnode/ethstats/stats"
)

const ID = "vipstats"

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("connected, upgrading", r)
		conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		go func() error {
			defer conn.Close()
			defer fmt.Println("disconnecting", conn)

			fmt.Println("upgraded", conn)
			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					return err
				}
				fmt.Println("received", string(msg), op, err)

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
						stats.NodePing{ID, time.Now()},
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
				fmt.Println("sent", string(out), op, err)
				//err = wsutil.WriteServerMessage(conn, op, msg)
				//if err != nil {
				//return err
				//}
			}
		}()
	}))
}
