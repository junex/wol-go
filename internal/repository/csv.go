package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// CSVFile CSV 文件操作基础结构
type CSVFile struct {
	filePath string
	mu       sync.RWMutex
}

// NewCSVFile 创建 CSV 文件操作实例
func NewCSVFile(filePath string) *CSVFile {
	return &CSVFile{
		filePath: filePath,
	}
}

// Read 读取 CSV 文件所有记录
func (f *CSVFile) Read() ([][]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	file, err := os.Open(f.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，返回空记录
			return [][]string{}, nil
		}
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取 CSV 失败: %w", err)
	}

	return records, nil
}

// Write 写入 CSV 文件（原子写入：临时文件 + 重命名）
func (f *CSVFile) Write(records [][]string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 确保目录存在
	dir := filepath.Dir(f.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建临时文件
	tempPath := f.filePath + ".tmp"
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}

	// 写入数据
	writer := csv.NewWriter(tempFile)
	if err := writer.WriteAll(records); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return fmt.Errorf("写入 CSV 失败: %w", err)
	}
	writer.Flush()
	tempFile.Close()

	// 原子性重命名
	if err := os.Rename(tempPath, f.filePath); err != nil {
		return fmt.Errorf("重命名文件失败: %w", err)
	}

	return nil
}

// Append 追加一条记录到 CSV 文件
func (f *CSVFile) Append(record []string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 确保目录存在
	dir := filepath.Dir(f.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 打开文件（追加模式）
	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 写入记录
	writer := csv.NewWriter(file)
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("写入记录失败: %w", err)
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

// Exists 检查文件是否存在
func (f *CSVFile) Exists() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, err := os.Stat(f.filePath)
	return err == nil
}
