package ws

import (
	ws_dto "converter/dto/ws"
	"log"
	"sync"
)

type Hub struct {
	mu         sync.RWMutex
	guests     map[uint]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	Notify     chan ws_dto.NotificationDto
}

func NewHub() *Hub {
	return &Hub{
		guests:     make(map[uint]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Notify:     make(chan ws_dto.NotificationDto, 100),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.guests[client.GuestID]; !ok {
				h.guests[client.GuestID] = make(map[*Client]bool)
			}
			h.guests[client.GuestID][client] = true
			h.mu.Unlock()
			log.Printf("Client with guest_id=%s connected (total tabs: %d)",
				client.GuestID, len(h.guests[client.GuestID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if guests, ok := h.guests[client.GuestID]; ok {
				if _, exists := guests[client]; exists {
					delete(guests, client)
					close(client.Send)
				}
				if len(guests) == 0 {
					delete(h.guests, client.GuestID)
				}
			}
			h.mu.Unlock()
			log.Printf("Client with guest_id=%s disconnected", client.GuestID)

		case notify := <-h.Notify:
			h.send(&notify)
		}
	}
}

func (h *Hub) send(notificationDto *ws_dto.NotificationDto) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	guests, ok := h.guests[notificationDto.GuestID]
	if !ok {
		log.Printf("No active connections for guest_id=%s", notificationDto.GuestID)
		return
	}

	msg := ws_dto.MessageDto{
		Type:    ws_dto.Notification,
		Payload: notificationDto,
	}

	for client := range guests {
		select {
		case client.Send <- msg:
		default:
			// Если канал клиента переполнен, пропускаем сообщение.
		}
	}
}

func (h *Hub) countGuest() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.guests)
}

func (h *Hub) totalConnections() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	total := 0
	for _, clients := range h.guests {
		total += len(clients)
	}
	return total
}
