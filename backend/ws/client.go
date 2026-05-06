package ws

import (
	ws_dto "converter/dto/ws"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	GuestID uint
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan ws_dto.MessageDto
}

func ServeWs(hub *Hub, guestID uint, w http.ResponseWriter, r *http.Request) {
	if guestID <= 0 {
		log.Println("error: not found guest_id")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	client := &Client{
		GuestID: guestID,
		Hub:     hub,
		Conn:    conn,
		Send:    make(chan ws_dto.MessageDto, 256),
	}

	hub.register <- client

	client.Send <- ws_dto.MessageDto{
		Type: ws_dto.RegisterAck,
		Payload: ws_dto.RegisterAckDto{
			GuestID: guestID,
			Status:  "ok",
		},
	}

	go client.writePump()
	go client.readPump()
}

// readPump читает сообщения от клиента (принимает только pong).
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("marshal error for guest %d: %v", c.GuestID, err)
				continue
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

			// Если в канале накопились ещё сообщения, отправляем каждое отдельно
			n := len(c.Send)
			for i := 0; i < n; i++ {
				nextMsg := <-c.Send
				nextData, err := json.Marshal(nextMsg)
				if err != nil {
					log.Printf("marshal error: %v", err)
					continue
				}
				if err := c.Conn.WriteMessage(websocket.TextMessage, nextData); err != nil {
					return
				}
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
