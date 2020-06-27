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

var conns = make(map[string]*websocket.Conn)

type message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

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

	go handleConn(conn, conns)
}

func handleConn(conn *websocket.Conn, conns map[string]*websocket.Conn) {
	for {
		msg := message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			conn.Close()
			delete(conns, msg.From)
			break
		}

		if _, ok := conns[msg.From]; !ok {
			conns[msg.From] = conn
		}

		if msg.To != "" {
			if out, ok := conns[msg.To]; ok {
				if err := out.WriteJSON(&msg); err != nil {
					log.Println(err)
					out.Close()
					delete(conns, msg.To)
				}
			}
			if err := conn.WriteJSON(&msg); err != nil {
				log.Println(err)
				conn.Close()
				delete(conns, msg.From)
			}
		} else {
			for k, v := range conns {
				if err := v.WriteJSON(&msg); err != nil {
					log.Println(err)
					v.Close()
					delete(conns, k)
					continue
				}
			}
		}
	}
}
