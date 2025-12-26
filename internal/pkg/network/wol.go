package network

import (
	"fmt"
	"net"
)

// SendMagicPacket 发送 WOL 魔术包
func SendMagicPacket(macAddr string, broadcastAddr string) error {
	// 解析 MAC 地址
	mac, err := net.ParseMAC(macAddr)
	if err != nil {
		return fmt.Errorf("无效的 MAC 地址: %w", err)
	}

	// 构造魔术包: 6×0xFF + 16×MAC地址
	packet := make([]byte, 102)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 6; i < 102; i += 6 {
		copy(packet[i:], mac)
	}

	// UDP 广播
	conn, err := net.Dial("udp", broadcastAddr+":9")
	if err != nil {
		return fmt.Errorf("创建 UDP 连接失败: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("发送魔术包失败: %w", err)
	}

	return nil
}
