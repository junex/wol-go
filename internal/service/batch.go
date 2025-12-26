package service

import (
	"log"
	"sync"

	"github.com/junex/wol-go/internal/model"
)

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	Success       bool              `json:"success"`
	Total         int               `json:"total"`
	SuccessCount  int               `json:"success_count"`
	FailureCount  int               `json:"failure_count"`
	Results       map[string]bool   `json:"results"`       // MAC -> 操作是否成功
	ErrorMessages map[string]string `json:"error_messages"` // MAC -> 错误消息
}

// BatchService 批量操作服务
type BatchService struct {
	wolService    *WOLService
	statusService *StatusService
	maxConcurrent int // 最大并发数
}

// NewBatchService 创建批量操作服务
func NewBatchService(wolService *WOLService, statusService *StatusService) *BatchService {
	return &BatchService{
		wolService:    wolService,
		statusService: statusService,
		maxConcurrent: 10, // 默认最多 10 个并发
	}
}

// BatchWake 批量唤醒设备
func (s *BatchService) BatchWake(macAddresses []string) BatchOperationResult {
	result := BatchOperationResult{
		Total:         len(macAddresses),
		Results:       make(map[string]bool),
		ErrorMessages: make(map[string]string),
	}

	// 使用工作池控制并发
	sem := make(chan struct{}, s.maxConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, mac := range macAddresses {
		wg.Add(1)
		go func(macAddr string) {
			defer wg.Done()
			sem <- struct{}{}        // 获取信号量
			defer func() { <-sem }() // 释放信号量

			err := s.wolService.WakeByMAC(macAddr)
			mu.Lock()
			if err != nil {
				result.Results[macAddr] = false
				result.ErrorMessages[macAddr] = err.Error()
				result.FailureCount++
			} else {
				result.Results[macAddr] = true
				result.SuccessCount++
			}
			mu.Unlock()
		}(mac)
	}
	wg.Wait()

	result.Success = result.SuccessCount > 0
	log.Printf("[BatchService] BatchWake completed: %d/%d succeeded",
		result.SuccessCount, result.Total)

	return result
}

// BatchSleep 批量关机
func (s *BatchService) BatchSleep(computers []model.Computer) BatchOperationResult {
	result := BatchOperationResult{
		Total:         len(computers),
		Results:       make(map[string]bool),
		ErrorMessages: make(map[string]string),
	}

	// 使用工作池控制并发
	sem := make(chan struct{}, s.maxConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, computer := range computers {
		wg.Add(1)
		go func(comp model.Computer) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			err := s.wolService.Sleep(comp)
			mu.Lock()
			if err != nil {
				result.Results[comp.MACAddr] = false
				result.ErrorMessages[comp.MACAddr] = err.Error()
				result.FailureCount++
			} else {
				result.Results[comp.MACAddr] = true
				result.SuccessCount++
			}
			mu.Unlock()
		}(computer)
	}
	wg.Wait()

	result.Success = result.SuccessCount > 0
	log.Printf("[BatchService] BatchSleep completed: %d/%d succeeded",
		result.SuccessCount, result.Total)

	return result
}

// BatchCheckStatus 批量检查设备状态
func (s *BatchService) BatchCheckStatus(computers []model.Computer) map[string]bool {
	log.Printf("[BatchService] BatchCheckStatus: checking %d computers", len(computers))
	return s.statusService.CheckAll(computers)
}

// SetMaxConcurrent 设置最大并发数
func (s *BatchService) SetMaxConcurrent(max int) {
	if max > 0 {
		s.maxConcurrent = max
		log.Printf("[BatchService] Max concurrent set to %d", s.maxConcurrent)
	}
}

// GetMaxConcurrent 获取最大并发数
func (s *BatchService) GetMaxConcurrent() int {
	return s.maxConcurrent
}
