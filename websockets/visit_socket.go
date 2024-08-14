package websockets

var clients = make(map[string]*WebsocketConnection) // connected clients

type Visit struct {
	Country string `json:"country"`
	Referer string `json:"referer"`
}

func SendLog(productID string, visit Visit) {
	if client, ok := clients[productID]; ok {
		client.SendData(visit)
	}
}

func OpenLogConnection(conn *WebsocketConnection, productID string) error {
	clients[productID] = conn
	return nil
}
