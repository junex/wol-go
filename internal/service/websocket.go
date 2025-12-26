package service

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSMessage WebSocket 消息结构
type WSMessage struct {
	Type    string      `json:"type"`    // 消息类型: status, computer_added, computer_deleted, computer_updated, etc.
	Payload interface{} `json:"payload"` // 消息数据
}

// WSClient WebSocket 客户端
type WSClient struct {
	Conn *websocket.Conn
	Send chan WSMessage
	Hub  *WSHub
}

// NewWSClient 创建新的 WebSocket 客户端
func NewWSClient(conn *websocket.Conn, hub *WSHub) *WSClient {
	return &WSClient{
		Conn: conn,
		Send: make(chan WSMessage, 256),
		Hub:  hub,
	}
}

// WSHub WebSocket 连接管理器
type WSHub struct {
	// 注册的客户端
	Clients map[*WSClient]bool

	// 注册请求
	Register chan *WSClient

	// 注销请求
	Unregister chan *WSClient

	// 广播消息
	Broadcast chan WSMessage

	// 互斥锁
	mu sync.RWMutex
}

// NewWSHub 创建新的 WebSocket Hub
func NewWSHub() *WSHub {
	return &WSHub{
		Clients:    make(map[*WSClient]bool),
		Register:   make(chan *WSClient),
		Unregister: make(chan *WSClient),
		Broadcast:  make(chan WSMessage),
	}
}

// Run 启动 Hub 事件循环
func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Printf("[WSHub] Client registered. Total clients: %d", len(h.Clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("[WSHub] Client unregistered. Total clients: %d", len(h.Clients))

		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					// 发送缓冲区满，关闭客户端
					log.Printf("[WSHub] Client send buffer full, closing connection")
					h.Unregister <- client
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastStatus 广播设备状态更新
func (h *WSHub) BroadcastStatus(macAddress string, online bool) {
	message := WSMessage{
		Type: "status",
		Payload: map[string]interface{}{
			"mac_address": macAddress,
			"online":      online,
		},
	}
	h.Broadcast <- message
}

// BroadcastComputerAdded 广播设备添加事件
func (h *WSHub) BroadcastComputerAdded(computer interface{}) {
	message := WSMessage{
		Type:    "computer_added",
		Payload: computer,
	}
	h.Broadcast <- message
}

// BroadcastComputerUpdated 广播设备更新事件
func (h *WSHub) BroadcastComputerUpdated(computer interface{}) {
	message := WSMessage{
		Type:    "computer_updated",
		Payload: computer,
	}
	h.Broadcast <- message
}

// BroadcastComputerDeleted 广播设备删除事件
func (h *WSHub) BroadcastComputerDeleted(macAddress string) {
	message := WSMessage{
		Type: "computer_deleted",
		Payload: map[string]interface{}{
			"mac_address": macAddress,
		},
	}
	h.Broadcast <- message
}

// BroadcastCronAdded 广播 Cron 任务添加事件
func (h *WSHub) BroadcastCronAdded(macAddress string, cronType string, schedule string) {
	message := WSMessage{
		Type: "cron_added",
		Payload: map[string]interface{}{
			"mac_address": macAddress,
			"type":        cronType,
			"schedule":    schedule,
		},
	}
	h.Broadcast <- message
}

// BroadcastCronDeleted 广播 Cron 任务删除事件
func (h *WSHub) BroadcastCronDeleted(macAddress string, cronType string) {
	message := WSMessage{
		Type: "cron_deleted",
		Payload: map[string]interface{}{
			"mac_address": macAddress,
			"type":        cronType,
		},
	}
	h.Broadcast <- message
}

// ClientCount 获取当前连接的客户端数量
func (h *WSHub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Clients)
}

// ReadPump 从 WebSocket 连接读取消息
func (c *WSClient) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// 设置读取限制
	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// 配置 Pong 处理器
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WSClient] Read error: %v", err)
			}
			break
		}

		// 处理接收到的消息（如果需要）
		log.Printf("[WSClient] Received message: %s", string(message))
	}
}

// WritePump 向 WebSocket 连接写入消息
func (c *WSClient) WritePump() {
	// 配置 Ping 定时器 (每 54 秒发送一次 Ping，超时时间 60 秒)
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 序列化消息
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("[WSClient] Failed to marshal message: %v", err)
				continue
			}

			// 发送消息
			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("[WSClient] Write error: %v", err)
				return
			}

		case <-ticker.C:
			// 发送 Ping
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
