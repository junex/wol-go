package service

import (
	"sync"

	"github.com/junex/wol-go/internal/model"
	"github.com/junex/wol-go/internal/pkg/network"
)

// StatusService 状态检测服务
type StatusService struct {
	pingTimeout int
	arpTimeout  int
	tcpTimeout  int
}

// NewStatusService 创建状态检测服务
func NewStatusService(pingTimeout, arpTimeout, tcpTimeout int) *StatusService {
	return &StatusService{
		pingTimeout: pingTimeout,
		arpTimeout:  arpTimeout,
		tcpTimeout:  tcpTimeout,
	}
}

// CheckOne 检查单个设备状态
func (s *StatusService) CheckOne(computer model.Computer) bool {
	return network.ParseTestTypeAndCheck(
		computer.IPAddress,
		computer.TestType,
		s.pingTimeout,
		s.arpTimeout,
		s.tcpTimeout,
	)
}

// CheckAll 检查所有设备状态（并发）
func (s *StatusService) CheckAll(computers []model.Computer) map[string]bool {
	results := make(map[string]bool)
	var mu sync.Mutex

	// 工作池，最多 10 个并发
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup

	for _, computer := range computers {
		wg.Add(1)
		go func(comp model.Computer) {
			defer wg.Done()

			// 获取信号量
			sem <- struct{}{}
			defer func() { <-sem }()

			// 检查状态
			awake := s.CheckOne(comp)

			// 保存结果
			mu.Lock()
			results[comp.MACAddr] = awake
			mu.Unlock()
		}(computer)
	}

	wg.Wait()
	return results
}

// CheckByIP 根据 IP 和测试类型检查
func (s *StatusService) CheckByIP(ipAddress, testType string) bool {
	return network.ParseTestTypeAndCheck(
		ipAddress,
		testType,
		s.pingTimeout,
		s.arpTimeout,
		s.tcpTimeout,
	)
}
