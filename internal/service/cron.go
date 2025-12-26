package service

import (
	"fmt"
	"strings"

	"github.com/junex/wol-go/internal/pkg/validator"
	"github.com/junex/wol-go/internal/repository"
)

// CronService Cron 任务管理服务
type CronService struct {
	repo *repository.CronRepository
}

// NewCronService 创建 Cron 服务
func NewCronService(repo *repository.CronRepository) *CronService {
	return &CronService{repo: repo}
}

// AddWakeCron 添加唤醒定时任务
func (s *CronService) AddWakeCron(macAddr, schedule string) error {
	// 验证 Cron 表达式
	if err := validator.ValidateCron(schedule); err != nil {
		return fmt.Errorf("无效的 Cron 表达式: %w", err)
	}

	// 创建 Cron 条目
	entry := repository.CronEntry{
		Schedule: schedule,
		User:     "root",
		MACAddr:  macAddr,
		Type:     "wake",
	}

	// 添加到文件
	return s.repo.Add(entry)
}

// AddSleepCron 添加关机定时任务
func (s *CronService) AddSleepCron(macAddr, schedule string) error {
	// 验证 Cron 表达式
	if err := validator.ValidateCron(schedule); err != nil {
		return fmt.Errorf("无效的 Cron 表达式: %w", err)
	}

	// 反转 MAC 地址（用于 Sleep on LAN）
	reversedMAC := s.reverseMAC(macAddr)

	// 创建 Cron 条目
	entry := repository.CronEntry{
		Schedule: schedule,
		User:     "root",
		MACAddr:  reversedMAC,
		Type:     "sleep",
	}

	// 添加到文件
	return s.repo.Add(entry)
}

// DeleteWakeCron 删除唤醒定时任务
func (s *CronService) DeleteWakeCron(macAddr string) error {
	return s.repo.Delete(macAddr)
}

// DeleteSleepCron 删除关机定时任务
func (s *CronService) DeleteSleepCron(macAddr string) error {
	// 反转 MAC 地址
	reversedMAC := s.reverseMAC(macAddr)
	return s.repo.Delete(reversedMAC)
}

// GetCrons 获取设备的所有定时任务
func (s *CronService) GetCrons(macAddr string) (map[string]string, error) {
	result := make(map[string]string)

	// 获取唤醒任务
	wakeEntries, err := s.repo.GetByMAC(macAddr)
	if err != nil {
		return nil, err
	}
	for _, entry := range wakeEntries {
		result["wake"] = entry.Schedule
	}

	// 获取关机任务
	sleepEntries, err := s.repo.GetByReversedMAC(macAddr)
	if err != nil {
		return nil, err
	}
	for _, entry := range sleepEntries {
		result["sleep"] = entry.Schedule
	}

	return result, nil
}

// reverseMAC 反转 MAC 地址
func (s *CronService) reverseMAC(macAddr string) string {
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
