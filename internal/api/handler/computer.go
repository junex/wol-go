package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/junex/wol-go/internal/model"
	"github.com/junex/wol-go/internal/service"
)

// ComputerHandler 设备管理处理器
type ComputerHandler struct {
	computerService *service.ComputerService
	statusService   *service.StatusService
	wolService      *service.WOLService
}

// NewComputerHandler 创建设备管理处理器
func NewComputerHandler(
	compSvc *service.ComputerService,
	statusSvc *service.StatusService,
	wolSvc *service.WOLService,
) *ComputerHandler {
	return &ComputerHandler{
		computerService: compSvc,
		statusService:   statusSvc,
		wolService:      wolSvc,
	}
}

// ListComputers 获取设备列表
func (h *ComputerHandler) ListComputers(w http.ResponseWriter, r *http.Request) {
	computers, err := h.computerService.GetAll()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    computers,
	})
}

// AddComputer 添加设备
func (h *ComputerHandler) AddComputer(w http.ResponseWriter, r *http.Request) {
	var computer model.Computer
	if err := json.NewDecoder(r.Body).Decode(&computer); err != nil {
		h.writeError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if err := h.computerService.Add(computer); err != nil {
		h.writeError(w, http.StatusConflict, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "设备添加成功",
		"data":    computer,
	})
}

// UpdateComputer 更新设备
func (h *ComputerHandler) UpdateComputer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["mac"]

	var computer model.Computer
	if err := json.NewDecoder(r.Body).Decode(&computer); err != nil {
		h.writeError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 确保 MAC 地址匹配
	computer.MACAddr = macAddr

	if err := h.computerService.Update(computer); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "设备更新成功",
		"data":    computer,
	})
}

// DeleteComputer 删除设备
func (h *ComputerHandler) DeleteComputer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["mac"]

	if err := h.computerService.Delete(macAddr); err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "设备删除成功",
	})
}

// GetStatus 获取设备状态
func (h *ComputerHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["mac"]

	computer, err := h.computerService.GetByMAC(macAddr)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	awake := h.statusService.CheckOne(*computer)

	status := "asleep"
	if awake {
		status = "awake"
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(status))
}

// WakeComputer 唤醒设备
func (h *ComputerHandler) WakeComputer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["mac"]

	computer, err := h.computerService.GetByMAC(macAddr)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := h.wolService.Wake(*computer); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "WOL 魔术包已发送到 " + computer.Name,
	})
}

// SleepComputer 关机设备
func (h *ComputerHandler) SleepComputer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["mac"]

	computer, err := h.computerService.GetByMAC(macAddr)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := h.wolService.Sleep(*computer); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "关机指令已发送到 " + computer.Name,
	})
}

// writeJSON 写入 JSON 响应
func (h *ComputerHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError 写入错误响应
func (h *ComputerHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"message": message,
		},
	})
}
