package service

import (
	"fmt"
	"strings"

	"github.com/junex/wol-go/internal/model"
	"github.com/junex/wol-go/internal/repository"
	"github.com/junex/wol-go/internal/pkg/validator"
)

// ComputerService 设备管理服务
type ComputerService struct {
	repo *repository.ComputerRepository
}

// NewComputerService 创建设备管理服务
func NewComputerService(repo *repository.ComputerRepository) *ComputerService {
	return &ComputerService{repo: repo}
}

// Add 添加设备
func (s *ComputerService) Add(computer model.Computer) error {
	// 验证输入
	if err := s.validateComputer(computer); err != nil {
		return err
	}

	// 检查是否已存在
	if existing, _ := s.repo.GetByName(computer.Name); existing != nil {
		return fmt.Errorf("设备名称 %s 已存在", computer.Name)
	}

	if existing, _ := s.repo.GetByMAC(computer.MACAddr); existing != nil {
		return fmt.Errorf("MAC 地址 %s 已存在", computer.MACAddr)
	}

	// 添加到仓库
	return s.repo.Add(computer)
}

// Update 更新设备
func (s *ComputerService) Update(computer model.Computer) error {
	// 验证输入
	if err := s.validateComputer(computer); err != nil {
		return err
	}

	// 更新
	return s.repo.Update(computer)
}

// Delete 删除设备
func (s *ComputerService) Delete(macAddr string) error {
	return s.repo.Delete(macAddr)
}

// GetAll 获取所有设备
func (s *ComputerService) GetAll() ([]model.Computer, error) {
	return s.repo.GetAll()
}

// GetByMAC 根据 MAC 地址获取设备
func (s *ComputerService) GetByMAC(macAddr string) (*model.Computer, error) {
	return s.repo.GetByMAC(macAddr)
}

// GetByName 根据名称获取设备
func (s *ComputerService) GetByName(name string) (*model.Computer, error) {
	return s.repo.GetByName(name)
}

// validateComputer 验证设备信息
func (s *ComputerService) validateComputer(computer model.Computer) error {
	if computer.Name == "" {
		return fmt.Errorf("设备名称不能为空")
	}

	if err := validator.ValidateMAC(computer.MACAddr); err != nil {
		return err
	}

	if err := validator.ValidateIP(computer.IPAddress); err != nil {
		return err
	}

	if err := validator.ValidateTestType(computer.TestType); err != nil {
		return err
	}

	// MAC 地址统一转大写
	computer.MACAddr = strings.ToUpper(computer.MACAddr)

	// TestType 统一转小写
	computer.TestType = strings.ToLower(computer.TestType)

	return nil
}
