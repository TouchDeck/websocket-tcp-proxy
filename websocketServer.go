package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type websocketClient struct {
	conn     *websocket.Conn
	serv     *websocketServer
	remoteIp string
}

type websocketServer struct {
	address           string
	onClientConnected func(c *websocketClient)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		// TODO: Actually check origin
		return true
	},
}

func (c *websocketClient) close() {
	c.conn.Close()
}

func (c *websocketClient) readMessage() (string, error) {
	_, msg, err := c.conn.ReadMessage()
	return string(msg), err
}

func (s *websocketServer) handleClient(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Could not upgrade HTTP request:", err)
		return
	}

	newClient := &websocketClient{
		conn:     conn,
		serv:     s,
		remoteIp: getRemoteIp(r.RemoteAddr),
	}
	s.onClientConnected(newClient)
}

func (s *websocketServer) listen(address string) {
	log.Println("Starting HTTP server on:", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalln("Error starting websocket server:", err)
	}
}

func newWebsocketServer(path string) *websocketServer {
	s := &websocketServer{
		onClientConnected: func(c *websocketClient) {},
	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		s.handleClient(w, r)
	})

	return s
}