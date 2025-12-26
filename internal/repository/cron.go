package repository

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// CronRepository Cron 任务数据访问层
type CronRepository struct {
	filePath string
	mu       sync.RWMutex
}

// NewCronRepository 创建 Cron 仓库实例
func NewCronRepository(filePath string) *CronRepository {
	return &CronRepository{
		filePath: filePath,
	}
}

// CronEntry Cron 任务条目
type CronEntry struct {
	Schedule string // Cron 表达式
	User     string // 用户
	MACAddr  string // MAC 地址
	Type     string // "wake" or "sleep"
}

// ReadAll 读取所有 Cron 任务
func (r *CronRepository) ReadAll() ([]CronEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，返回空列表
			return []CronEntry{}, nil
		}
		return nil, fmt.Errorf("打开 cron 文件失败: %w", err)
	}
	defer file.Close()

	entries := []CronEntry{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 cron 行
		entry, err := r.parseCronLine(line)
		if err != nil {
			// 跳过无效行，但记录日志
			continue
		}

		entries = append(entries, *entry)
	}

	return entries, scanner.Err()
}

// parseCronLine 解析一行 cron 配置
// 格式: "分 时 日 月 周 用户 命令"
// 示例: "0 8 * * * root /usr/local/bin/wakeonlan 00:11:22:33:44:55"
func (r *CronRepository) parseCronLine(line string) (*CronEntry, error) {
	fields := strings.Fields(line)
	if len(fields) < 7 {
		return nil, fmt.Errorf("无效的 cron 行: %s", line)
	}

	// 前 5 个字段是 cron 表达式
	schedule := strings.Join(fields[0:5], " ")

	// 第 6 个字段是用户
	user := fields[5]

	// 剩余字段是命令
	command := strings.Join(fields[6:], " ")

	// 从命令中提取 MAC 地址
	// 命令格式: "/usr/local/bin/wakeonlan XX:XX:XX:XX:XX:XX"
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return nil, fmt.Errorf("无效的命令: %s", command)
	}

	macAddr := parts[len(parts)-1]

	// 判断类型：如果是反转的 MAC，则是 sleep 任务
	// 原始 MAC: 00:11:22:33:44:55
	// 反转 MAC: 55:44:33:22:11:00
	entry := &CronEntry{
		Schedule: schedule,
		User:     user,
		MACAddr:  macAddr,
		Type:     "wake", // 默认为 wake
	}

	return entry, nil
}

// Add 添加 Cron 任务
func (r *CronRepository) Add(entry CronEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 构造 cron 行
	line := r.buildCronLine(entry)

	// 追加到文件
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开 cron 文件失败: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(line + "\n"); err != nil {
		return fmt.Errorf("写入 cron 文件失败: %w", err)
	}

	return nil
}

// Delete 删除指定 MAC 地址的 Cron 任务
func (r *CronRepository) Delete(macAddr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 读取所有行
	lines, err := r.readAllLines()
	if err != nil {
		return err
	}

	// 过滤掉要删除的行
	newLines := []string{}
	deleted := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 保留注释和空行
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			newLines = append(newLines, line)
			continue
		}

		// 解析行
		entry, err := r.parseCronLine(trimmed)
		if err != nil {
			// 解析失败，保留该行
			newLines = append(newLines, line)
			continue
		}

		// 如果 MAC 地址不匹配，保留该行
		if !strings.EqualFold(entry.MACAddr, macAddr) {
			newLines = append(newLines, line)
		} else {
			deleted = true
		}
	}

	if !deleted {
		return fmt.Errorf("未找到 MAC 地址为 %s 的 cron 任务", macAddr)
	}

	// 写回文件
	return r.writeAllLines(newLines)
}

// buildCronLine 构造 cron 行
func (r *CronRepository) buildCronLine(entry CronEntry) string {
	command := fmt.Sprintf("/usr/local/bin/wakeonlan %s", entry.MACAddr)
	return fmt.Sprintf("%s %s %s", entry.Schedule, entry.User, command)
}

// readAllLines 读取文件所有行
func (r *CronRepository) readAllLines() ([]string, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// writeAllLines 写入所有行
func (r *CronRepository) writeAllLines(lines []string) error {
	// 确保目录存在
	dir := r.filePath[:strings.LastIndex(r.filePath, "/")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 创建临时文件
	tempPath := r.filePath + ".tmp"
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// 写入所有行
	for _, line := range lines {
		if _, err := tempFile.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	// 原子性重命名
	if err := os.Rename(tempPath, r.filePath); err != nil {
		return err
	}

	return nil
}

// GetByMAC 根据原始 MAC 地址获取 Cron 任务
func (r *CronRepository) GetByMAC(macAddr string) ([]CronEntry, error) {
	entries, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	result := []CronEntry{}
	for _, entry := range entries {
		if strings.EqualFold(entry.MACAddr, macAddr) {
			result = append(result, entry)
		}
	}

	return result, nil
}

// GetByReversedMAC 根据反转的 MAC 地址获取 Cron 任务（用于 SoL）
func (r *CronRepository) GetByReversedMAC(macAddr string) ([]CronEntry, error) {
	entries, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	reversedMAC := reverseMAC(macAddr)

	result := []CronEntry{}
	for _, entry := range entries {
		if strings.EqualFold(entry.MACAddr, reversedMAC) {
			result = append(result, entry)
		}
	}

	return result, nil
}

// reverseMAC 反转 MAC 地址
// 例如: 00:11:22:33:44:55 -> 55:44:33:22:11:00
func reverseMAC(macAddr string) string {
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
