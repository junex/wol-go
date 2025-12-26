package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// AuthService 认证服务
type AuthService struct {
	username    string
	password    string
	enableLogin bool
	sessions    map[string]*Session
	mu          sync.RWMutex
}

// Session 会话
type Session struct {
	Token     string
	Username  string
	ExpiresAt time.Time
}

// NewAuthService 创建认证服务
func NewAuthService(username, password string, enableLogin bool) *AuthService {
	return &AuthService{
		username:    username,
		password:    password,
		enableLogin: enableLogin,
		sessions:    make(map[string]*Session),
	}
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (string, error) {
	if !s.enableLogin {
		return "", fmt.Errorf("登录功能未启用")
	}

	if username != s.username || password != s.password {
		return "", fmt.Errorf("用户名或密码错误")
	}

	// 生成随机 token
	token, err := s.generateToken()
	if err != nil {
		return "", fmt.Errorf("生成 token 失败: %w", err)
	}

	// 创建会话（24小时有效期）
	session := &Session{
		Token:     token,
		Username:  username,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.mu.Lock()
	s.sessions[token] = session
	s.mu.Unlock()

	return token, nil
}

// Logout 用户登出
func (s *AuthService) Logout(token string) error {
	if !s.enableLogin {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, token)
	return nil
}

// ValidateToken 验证 token
func (s *AuthService) ValidateToken(token string) (*Session, error) {
	if !s.enableLogin {
		// 登录未启用，返回一个默认会话
		return &Session{
			Token:     "default",
			Username:  "admin",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[token]
	if !exists {
		return nil, fmt.Errorf("无效的 token")
	}

	// 检查是否过期
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, token)
		return nil, fmt.Errorf("token 已过期")
	}

	return session, nil
}

// CleanupExpiredSessions 清理过期会话
func (s *AuthService) CleanupExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, token)
		}
	}
}

// generateToken 生成随机 token
func (s *AuthService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
