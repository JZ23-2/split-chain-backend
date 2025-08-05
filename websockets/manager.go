package websockets

import "github.com/gorilla/websocket"

var clients = make(map[string]*websocket.Conn)

func RegisterClient(userID string, conn *websocket.Conn) {
	clients[userID] = conn
}

func UnregisterClient(userID string) {
	delete(clients, userID)
}

func SendToUser(userID string, message any) error {
	if conn, ok := clients[userID]; ok {
		return conn.WriteJSON(message)
	}
	return nil
}
