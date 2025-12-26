package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/junex/wol-go/internal/model"
	"github.com/junex/wol-go/internal/pkg/network"
)

// WOLService WOL 控制服务
type WOLService struct {
	statusService *StatusService
}

// NewWOLService 创建 WOL 服务
func NewWOLService(statusService *StatusService) *WOLService {
	return &WOLService{
		statusService: statusService,
	}
}

// Wake 唤醒设备
func (s *WOLService) Wake(computer model.Computer) error {
	// 发送 WOL 魔术包
	err := network.SendMagicPacket(computer.MACAddr, "255.255.255.255")
	if err != nil {
		return fmt.Errorf("发送 WOL 魔术包失败: %w", err)
	}

	return nil
}

// WakeByMAC 通过 MAC 地址唤醒设备（简化版本）
func (s *WOLService) WakeByMAC(macAddr string) error {
	// 发送 WOL 魔术包
	err := network.SendMagicPacket(macAddr, "255.255.255.255")
	if err != nil {
		return fmt.Errorf("发送 WOL 魔术包失败: %w", err)
	}

	return nil
}

// Sleep 关机设备
func (s *WOLService) Sleep(computer model.Computer) error {
	// 首先检查设备是否在线
	if !s.statusService.CheckOne(computer) {
		return fmt.Errorf("设备 %s 不在线，无法发送关机指令", computer.Name)
	}

	// 尝试通过 HTTP 发送关机命令
	shutdownURL := fmt.Sprintf("http://%s:8009/shutdown", computer.IPAddress)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(shutdownURL)
	if err != nil {
		return fmt.Errorf("连接关机服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("关机服务返回错误状态码: %d", resp.StatusCode)
	}

	return nil
}

// WakeOrSleep 根据设备状态决定唤醒或关机
func (s *WOLService) WakeOrSleep(computer model.Computer) (string, error) {
	isAwake := s.statusService.CheckOne(computer)

	if isAwake {
		// 设备在线，发送关机指令
		err := s.Sleep(computer)
		if err != nil {
			return "sleep", err
		}
		return "sleep", nil
	} else {
		// 设备离线，发送唤醒包
		err := s.Wake(computer)
		if err != nil {
			return "wake", err
		}
		return "wake", nil
	}
}

// ReverseMAC 反转 MAC 地址（用于 SoL）
func ReverseMAC(macAddr string) string {
	parts := strings.Split(macAddr, ":")
	if len(parts) != 6 {
		return macAddr
	}

	// 反转数组
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	return strings.Join(parts, ":")
}
