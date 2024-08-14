package websockets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketConnection struct {
	conn *websocket.Conn
}

func (w *WebsocketConnection) SendErrorMesage(message string) error {
	return w.conn.WriteJSON(gin.H{"error": message})
}

func (w *WebsocketConnection) SendData(data interface{}) error {
	return w.conn.WriteJSON(data)
}

func (w *WebsocketConnection) Close() error {
	return w.conn.Close()
}

func (w *WebsocketConnection) ReadMessage() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *WebsocketConnection) LoopConnection() {
	for {
		_, _, err := w.ReadMessage()
		if err != nil {
			break
		}
	}
}

func UpgradeConnection(c *gin.Context) (*WebsocketConnection, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	return &WebsocketConnection{conn: conn}, nil
}
