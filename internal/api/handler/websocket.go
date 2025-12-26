package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/junex/wol-go/internal/service"
)

// WebSocketUpgrader WebSocket 升级器
var WebSocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有来源（可以根据需要限制）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler WebSocket 处理器
type WebSocketHandler struct {
	hub *service.WSHub
}

// NewWebSocketHandler 创建 WebSocket 处理器
func NewWebSocketHandler(hub *service.WSHub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接到 WebSocket
	conn, err := WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WebSocket] Failed to upgrade connection: %v", err)
		return
	}

	// 创建客户端
	client := service.NewWSClient(conn, h.hub)

	// 注册客户端
	h.hub.Register <- client

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()

	log.Printf("[WebSocket] New connection from %s", r.RemoteAddr)
}

// BroadcastStatus 广播状态更新（辅助方法）
func (h *WebSocketHandler) BroadcastStatus(macAddress string, online bool) {
	h.hub.BroadcastStatus(macAddress, online)
}

// BroadcastComputerAdded 广播设备添加
func (h *WebSocketHandler) BroadcastComputerAdded(computer interface{}) {
	h.hub.BroadcastComputerAdded(computer)
}

// BroadcastComputerUpdated 广播设备更新
func (h *WebSocketHandler) BroadcastComputerUpdated(computer interface{}) {
	h.hub.BroadcastComputerUpdated(computer)
}

// BroadcastComputerDeleted 广播设备删除
func (h *WebSocketHandler) BroadcastComputerDeleted(macAddress string) {
	h.hub.BroadcastComputerDeleted(macAddress)
}

// GetClientCount 获取连接数
func (h *WebSocketHandler) GetClientCount() int {
	return h.hub.ClientCount()
}

