package repository

import (
	"fmt"
	"strings"
	"sync"

	"github.com/junex/wol-go/internal/model"
)

// ComputerRepository 计算机数据访问层
type ComputerRepository struct {
	csv     *CSVFile
	cache   []model.Computer
	mu      sync.RWMutex
	loaded  bool
}

// NewComputerRepository 创建计算机仓库实例
func NewComputerRepository(filePath string) *ComputerRepository {
	return &ComputerRepository{
		csv: NewCSVFile(filePath),
	}
}

// LoadAll 从文件加载所有计算机到内存缓存
func (r *ComputerRepository) LoadAll() ([]model.Computer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 从 CSV 文件读取
	records, err := r.csv.Read()
	if err != nil {
		return nil, err
	}

	// 解析记录
	computers := make([]model.Computer, 0, len(records))
	for _, record := range records {
		if len(record) < 3 {
			continue // 跳过无效记录
		}

		computer := model.Computer{
			Name:      record[0],
			MACAddr:   record[1],
			IPAddress: record[2],
		}

		// TestType 是可选的，默认为 "icmp"
		if len(record) >= 4 && record[3] != "" {
			computer.TestType = record[3]
		} else {
			computer.TestType = "icmp"
		}

		computers = append(computers, computer)
	}

	// 更新缓存
	r.cache = computers
	r.loaded = true

	return computers, nil
}

// GetAll 获取所有计算机（从缓存）
func (r *ComputerRepository) GetAll() ([]model.Computer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.loaded {
		return nil, fmt.Errorf("数据未加载，请先调用 LoadAll()")
	}

	// 返回副本，避免外部修改
	result := make([]model.Computer, len(r.cache))
	copy(result, r.cache)
	return result, nil
}

// GetByMAC 根据 MAC 地址查找计算机
func (r *ComputerRepository) GetByMAC(macAddr string) (*model.Computer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.loaded {
		return nil, fmt.Errorf("数据未加载")
	}

	for i := range r.cache {
		if strings.EqualFold(r.cache[i].MACAddr, macAddr) {
			// 返回副本
			computer := r.cache[i]
			return &computer, nil
		}
	}

	return nil, fmt.Errorf("未找到 MAC 地址为 %s 的计算机", macAddr)
}

// GetByName 根据名称查找计算机
func (r *ComputerRepository) GetByName(name string) (*model.Computer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.loaded {
		return nil, fmt.Errorf("数据未加载")
	}

	for i := range r.cache {
		if r.cache[i].Name == name {
			// 返回副本
			computer := r.cache[i]
			return &computer, nil
		}
	}

	return nil, fmt.Errorf("未找到名称为 %s 的计算机", name)
}

// Add 添加计算机
func (r *ComputerRepository) Add(computer model.Computer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查是否已存在
	for _, c := range r.cache {
		if c.Name == computer.Name {
			return fmt.Errorf("计算机名称 %s 已存在", computer.Name)
		}
		if strings.EqualFold(c.MACAddr, computer.MACAddr) {
			return fmt.Errorf("MAC 地址 %s 已存在", computer.MACAddr)
		}
	}

	// 添加到缓存
	r.cache = append(r.cache, computer)

	// 持久化到文件
	if err := r.saveToFile(); err != nil {
		// 回滚
		r.cache = r.cache[:len(r.cache)-1]
		return err
	}

	return nil
}

// Update 更新计算机
func (r *ComputerRepository) Update(computer model.Computer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 查找并更新
	found := false
	for i := range r.cache {
		if strings.EqualFold(r.cache[i].MACAddr, computer.MACAddr) {
			// 检查名称是否与其他设备冲突
			for j, c := range r.cache {
				if j != i && c.Name == computer.Name {
					return fmt.Errorf("计算机名称 %s 已存在", computer.Name)
				}
			}

			// 更新（MAC 地址不可变）
			r.cache[i].Name = computer.Name
			r.cache[i].IPAddress = computer.IPAddress
			r.cache[i].TestType = computer.TestType
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("未找到 MAC 地址为 %s 的计算机", computer.MACAddr)
	}

	// 持久化到文件
	if err := r.saveToFile(); err != nil {
		return err
	}

	return nil
}

// Delete 删除计算机
func (r *ComputerRepository) Delete(macAddr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 查找并删除
	found := false
	newCache := make([]model.Computer, 0, len(r.cache)-1)
	for _, c := range r.cache {
		if !strings.EqualFold(c.MACAddr, macAddr) {
			newCache = append(newCache, c)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("未找到 MAC 地址为 %s 的计算机", macAddr)
	}

	// 更新缓存
	oldCache := r.cache
	r.cache = newCache

	// 持久化到文件
	if err := r.saveToFile(); err != nil {
		// 回滚
		r.cache = oldCache
		return err
	}

	return nil
}

// saveToFile 将当前缓存保存到 CSV 文件
func (r *ComputerRepository) saveToFile() error {
	records := make([][]string, len(r.cache))
	for i, computer := range r.cache {
		records[i] = []string{
			computer.Name,
			computer.MACAddr,
			computer.IPAddress,
			computer.TestType,
		}
	}

	return r.csv.Write(records)
}

// Exists 检查文件是否存在
func (r *ComputerRepository) Exists() bool {
	return r.csv.Exists()
}
