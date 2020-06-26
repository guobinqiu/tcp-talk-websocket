package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var conns = make(map[*websocket.Conn]*websocket.Conn)

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("0.0.0.0:7777", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conns[conn] = conn

	go handleConn(conn, conns)
}

func handleConn(conn *websocket.Conn, conns map[*websocket.Conn]*websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			conn.Close()
			delete(conns, conn)
			break
		}

		for out := range conns {
			if err := out.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
				out.Close()
				delete(conns, out)
				continue
			}
		}
	}
}
