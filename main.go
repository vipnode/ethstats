package main

import (
	"fmt"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

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
				//err = wsutil.WriteServerMessage(conn, op, msg)
				//if err != nil {
				//return err
				//}
			}
		}()
	}))
}
