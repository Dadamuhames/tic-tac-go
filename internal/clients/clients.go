package clients

import (
	"log"
	"maps"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id   int
	Name string
	Conn *websocket.Conn
	mu   sync.Mutex
}

type ConnectionPool struct {
	clients map[int]*Client
	mu      sync.RWMutex
	add     chan *Client
	remove  chan *Client
}

func NewConnectionPool() *ConnectionPool {
	cp := &ConnectionPool{
		clients: make(map[int]*Client),
		add:     make(chan *Client),
		remove:  make(chan *Client),
	}

	go cp.run()
	return cp
}

func (cp *ConnectionPool) run() {
	for {
		select {
		case client := <-cp.add:
			cp.mu.Lock()
			cp.clients[client.Id] = client
			cp.mu.Unlock()
			log.Printf("Client %d connected. Total: %d", client.Id, len(cp.clients))

		case client := <-cp.remove:
			cp.mu.Lock()
			if _, exists := cp.clients[client.Id]; exists {
				client.mu.Lock()
				if client.Conn != nil {
					client.Conn.Close()
					client.Conn = nil
				}
				client.mu.Unlock()
				log.Printf("Client %d disconnected. Total: %d", client.Id, len(cp.clients))
			}
			cp.mu.Unlock()
		}
	}
}

func (cp *ConnectionPool) AddClient(name string, conn *websocket.Conn) int {
	existing := cp.findByName(name)

	if existing != nil {
		existing.mu.Lock()
		defer existing.mu.Unlock()
		existing.Conn = conn
		return existing.Id
	}

	ids := slices.Collect(maps.Keys(cp.clients))

	newId := 1

	if len(ids) != 0 {
		newId = ids[len(ids)-1] + 1
	}

	client := Client{Id: newId, Name: name, Conn: conn}

	cp.add <- &client

	return newId
}

func (cp *ConnectionPool) findByName(name string) *Client {
	clients := cp.GetAll()

	for _, client := range clients {
		if client.Name == name {
			return client
		}
	}

	return nil
}

func (cp *ConnectionPool) RemoveClient(id int) {
	client := cp.GetClientById(id)
	cp.remove <- client
}

func (cp *ConnectionPool) GetClientById(id int) *Client {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.clients[id]
}

func (cp *ConnectionPool) GetAll() []*Client {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	clients := make([]*Client, 0, len(cp.clients))
	for _, client := range cp.clients {
		if client.Conn != nil {
			clients = append(clients, client)
		}
	}

	return clients
}

func (cp *ConnectionPool) BroadcastMessage(fromUserId int, message []byte) {
	cp.mu.RLock()

	clients := make([]*Client, 0, len(cp.clients))

	for _, client := range cp.clients {
		if client.Id != fromUserId {
			clients = append(clients, client)
		}
	}

	cp.mu.RUnlock()

	for _, client := range clients {
		client.SendMessage(message)
	}
}

func (c *Client) SendMessage(message []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Conn == nil {
		return websocket.ErrCloseSent
	}

	return c.Conn.WriteMessage(websocket.TextMessage, message)
}
