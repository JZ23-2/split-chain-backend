package handlers

import (
	"fmt"
	"net/http"

	"github.com/JZ23-2/splitbill-backend/websockets"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(c *gin.Context) {
	userID := c.Param("userId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Gagal upgrade:", err)
		return
	}

	fmt.Println("WebSocket connected:", userID)
	websockets.RegisterClient(userID, conn)
	defer func() {
		websockets.UnregisterClient(userID)
		conn.Close()
		fmt.Println("WebSocket disconnected:", userID)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
