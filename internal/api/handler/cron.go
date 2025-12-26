package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/junex/wol-go/internal/service"
)

// CronHandler Cron 任务处理器
type CronHandler struct {
	cronService *service.CronService
}

// NewCronHandler 创建 Cron 处理器
func NewCronHandler(cronService *service.CronService) *CronHandler {
	return &CronHandler{
		cronService: cronService,
	}
}

// GetCrons 获取设备的定时任务
func (h *CronHandler) GetCrons(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["mac"]

	crons, err := h.cronService.GetCrons(macAddr)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    crons,
	})
}

// AddWakeCron 添加唤醒定时任务
func (h *CronHandler) AddWakeCron(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["mac"]

	var req struct {
		Schedule string `json:"schedule"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if err := h.cronService.AddWakeCron(macAddr, req.Schedule); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "唤醒定时任务已添加",
		"data": map[string]interface{}{
			"mac":      macAddr,
			"schedule": req.Schedule,
			"type":     "wake",
		},
	})
}

// AddSleepCron 添加关机定时任务
func (h *CronHandler) AddSleepCron(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["mac"]

	var req struct {
		Schedule string `json:"schedule"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if err := h.cronService.AddSleepCron(macAddr, req.Schedule); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "关机定时任务已添加",
		"data": map[string]interface{}{
			"mac":      macAddr,
			"schedule": req.Schedule,
			"type":     "sleep",
		},
	})
}

// DeleteWakeCron 删除唤醒定时任务
func (h *CronHandler) DeleteWakeCron(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["mac"]

	if err := h.cronService.DeleteWakeCron(macAddr); err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "唤醒定时任务已删除",
	})
}

// DeleteSleepCron 删除关机定时任务
func (h *CronHandler) DeleteSleepCron(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	macAddr := vars["mac"]

	if err := h.cronService.DeleteSleepCron(macAddr); err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "关机定时任务已删除",
	})
}

// writeJSON 写入 JSON 响应
func (h *CronHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError 写入错误响应
func (h *CronHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"message": message,
		},
	})
}
