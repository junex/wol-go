package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/junex/wol-go/internal/model"
	"github.com/junex/wol-go/internal/service"
)

// BatchHandler 批量操作处理器
type BatchHandler struct {
	computerService *service.ComputerService
	batchService    *service.BatchService
}

// NewBatchHandler 创建批量操作处理器
func NewBatchHandler(
	computerService *service.ComputerService,
	batchService *service.BatchService,
) *BatchHandler {
	return &BatchHandler{
		computerService: computerService,
		batchService:    batchService,
	}
}

// BatchWakeRequest 批量唤醒请求
type BatchWakeRequest struct {
	MACAddresses []string `json:"mac_addresses"`
}

// BatchSleepRequest 批量关机请求
type BatchSleepRequest struct {
	MACAddresses []string `json:"mac_addresses"`
}

// BatchStatusRequest 批量状态检查请求
type BatchStatusRequest struct {
	MACAddresses []string `json:"mac_addresses"`
}

// BatchWake 批量唤醒设备
// POST /api/computers/batch/wake
func (h *BatchHandler) BatchWake(w http.ResponseWriter, r *http.Request) {
	// 解析请求
	var req BatchWakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证
	if len(req.MACAddresses) == 0 {
		http.Error(w, "MAC addresses list is empty", http.StatusBadRequest)
		return
	}

	// 执行批量唤醒
	result := h.batchService.BatchWake(req.MACAddresses)

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": result.Success,
		"data":    result,
	})
}

// BatchSleep 批量关机
// POST /api/computers/batch/sleep
func (h *BatchHandler) BatchSleep(w http.ResponseWriter, r *http.Request) {
	// 解析请求
	var req BatchSleepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证
	if len(req.MACAddresses) == 0 {
		http.Error(w, "MAC addresses list is empty", http.StatusBadRequest)
		return
	}

	// 获取设备信息
	computers := make([]model.Computer, 0, len(req.MACAddresses))
	for _, mac := range req.MACAddresses {
		computer, err := h.computerService.GetByMAC(mac)
		if err != nil {
			http.Error(w, "Computer not found: "+mac, http.StatusNotFound)
			return
		}
		computers = append(computers, *computer)
	}

	// 执行批量关机
	result := h.batchService.BatchSleep(computers)

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": result.Success,
		"data":    result,
	})
}

// BatchStatus 批量检查状态
// GET /api/computers/batch/status?mac=xx,xx,xx
func (h *BatchHandler) BatchStatus(w http.ResponseWriter, r *http.Request) {
	// 从查询参数获取 MAC 地址列表
	macList := r.URL.Query().Get("mac")
	if macList == "" {
		http.Error(w, "MAC addresses list is empty", http.StatusBadRequest)
		return
	}

	// 解析 MAC 地址列表
	var macAddresses []string
	// 简单的逗号分隔解析（可以使用更复杂的解析逻辑）
	for _, mac := range splitByComma(macList) {
		macAddresses = append(macAddresses, trimSpace(mac))
	}

	if len(macAddresses) == 0 {
		http.Error(w, "No valid MAC addresses", http.StatusBadRequest)
		return
	}

	// 获取设备信息
	computers := make([]model.Computer, 0, len(macAddresses))
	for _, mac := range macAddresses {
		computer, err := h.computerService.GetByMAC(mac)
		if err != nil {
			continue // 跳过找不到的设备
		}
		computers = append(computers, *computer)
	}

	if len(computers) == 0 {
		http.Error(w, "No valid computers found", http.StatusNotFound)
		return
	}

	// 执行批量状态检查
	results := h.batchService.BatchCheckStatus(computers)

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"total":   len(results),
			"results": results,
		},
	})
}

// 辅助函数
func splitByComma(s string) []string {
	var result []string
	start := 0
	for i, c := range s {
		if c == ',' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// RegisterRoutes 注册批量操作路由
func (h *BatchHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/computers/batch/wake", h.BatchWake).Methods("POST")
	router.HandleFunc("/api/computers/batch/sleep", h.BatchSleep).Methods("POST")
	router.HandleFunc("/api/computers/batch/status", h.BatchStatus).Methods("GET")
}
