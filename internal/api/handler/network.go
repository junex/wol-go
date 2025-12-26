package handler

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
)

// NetworkHandler 网络扫描处理器
type NetworkHandler struct {
	arpTimeout   int
	arpInterface string
}

// NewNetworkHandler 创建网络扫描处理器
func NewNetworkHandler(arpTimeout int, arpInterface string) *NetworkHandler {
	return &NetworkHandler{
		arpTimeout:   arpTimeout,
		arpInterface: arpInterface,
	}
}

// ARPScanResult ARP 扫描结果
type ARPScanResult struct {
	IP  string `json:"ip"`
	MAC string `json:"mac"`
}

// ARPScan ARP 网络扫描
func (h *NetworkHandler) ARPScan(w http.ResponseWriter, r *http.Request) {
	// 构造 arp-scan 命令
	args := []string{"-lqx", "-t", string(rune(h.arpTimeout))}

	if h.arpInterface != "" {
		args = append([]string{"-I", h.arpInterface}, args...)
	}

	cmd := exec.Command("arp-scan", args...)
	output, err := cmd.Output()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "ARP 扫描失败: "+err.Error())
		return
	}

	// 解析输出
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	devices := []ARPScanResult{}

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			devices = append(devices, ARPScanResult{
				IP:  parts[0],
				MAC: parts[1],
			})
		}
	}

	if len(devices) == 0 {
		h.writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "未发现新设备",
		})
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    devices,
	})
}

// writeJSON 写入 JSON 响应
func (h *NetworkHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError 写入错误响应
func (h *NetworkHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"message": message,
		},
	})
}
