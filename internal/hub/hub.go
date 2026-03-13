package hub

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Client struct {
	ID      string
	Channel string
	Send    chan []byte
}

type Message struct {
	Channel   string          `json:"-"`
	EventType string          `json:"eventType"`
	Message   string          `json:"msg"`
	Extra     json.RawMessage `json:"extra,omitempty"`
	Timestamp string          `json:"timestamp"`
}

type Hub struct {
	channels  map[string]map[*Client]bool
	register  chan *Client
	broadcast chan *Message
	mu        sync.RWMutex
}

func New() *Hub {
	return &Hub{
		channels:  make(map[string]map[*Client]bool),
		register:  make(chan *Client),
		broadcast: make(chan *Message, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.channels[client.Channel] == nil {
				h.channels[client.Channel] = make(map[*Client]bool)
			}
			h.channels[client.Channel][client] = true
			h.mu.Unlock()
			log.Printf("Client connected to channel: %s", client.Channel)

		case msg := <-h.broadcast:
			h.publishToChannel(msg)
		}
	}
}

func (h *Hub) Subscribe(channel string) *Client {
	client := &Client{
		ID:      generateID(),
		Channel: channel,
		Send:    make(chan []byte, 32),
	}
	h.register <- client
	return client
}

func (h *Hub) Unsubscribe(client *Client) {
	h.mu.Lock()
	if clients, ok := h.channels[client.Channel]; ok {
		if _, exists := clients[client]; exists {
			close(client.Send)
			delete(clients, client)
		}
		if len(clients) == 0 {
			delete(h.channels, client.Channel)
		}
	}
	h.mu.Unlock()
	log.Printf("Client disconnected from channel: %s", client.Channel)
}

func (h *Hub) Publish(channel string, msg *Message) {
	msg.Channel = channel
	h.broadcast <- msg
}

func (h *Hub) Broadcast(msg *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for channel := range h.channels {
		msgCopy := *msg
		msgCopy.Channel = channel
		h.publishToChannel(&msgCopy)
	}
}

func (h *Hub) publishToChannel(msg *Message) {
	h.mu.RLock()
	clients := make([]*Client, 0)
	if channelClients, ok := h.channels[msg.Channel]; ok {
		for client := range channelClients {
			clients = append(clients, client)
		}
	}
	h.mu.RUnlock()

	if len(clients) == 0 {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for _, client := range clients {
		select {
		case client.Send <- data:
		default:
			go h.Unsubscribe(client)
		}
	}
}

func (h *Hub) GetStats() (channels, clients int) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	channels = len(h.channels)
	for _, c := range h.channels {
		clients += len(c)
	}
	return
}

func (h *Hub) GetChannelClients(channel string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if c, ok := h.channels[channel]; ok {
		return len(c)
	}
	return 0
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
