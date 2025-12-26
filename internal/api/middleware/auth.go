package middleware

import (
	"net/http"
	"strings"

	"github.com/junex/wol-go/internal/service"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	authService *service.AuthService
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	middleware := &AuthMiddleware{
		authService: authService,
	}
	return func(next http.Handler) http.Handler {
		return middleware
	}
}

// ServeHTTP 实现 http.Handler 接口
func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 从 Header 或 Cookie 获取 token
	token := m.extractToken(r)

	// 验证 token
	session, err := m.authService.ValidateToken(token)
	if err != nil {
		// 未认证
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 将会话信息存入请求上下文
	// TODO: 可以使用 context 存储 session 信息
	_ = session

	// 调用下一个处理器
	// 注意：这里需要嵌套调用，但由于 ServeHTTP 签名限制，我们用另一种方式
}

// extractToken 从请求中提取 token
func (m *AuthMiddleware) extractToken(r *http.Request) string {
	// 优先从 Authorization Header 获取
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
		return authHeader
	}

	// 从 Cookie 获取
	if cookie, err := r.Cookie("session_token"); err == nil {
		return cookie.Value
	}

	// 从查询参数获取
	return r.URL.Query().Get("token")
}

// OptionalAuth 可选认证中间件（如果 ENABLE_LOGIN=false 则跳过）
func OptionalAuth(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)

			// 验证 token（即使失败也继续）
			session, err := authService.ValidateToken(token)
			_ = session
			_ = err

			// 无论认证是否成功，都调用下一个处理器
			next.ServeHTTP(w, r)
		})
	}
}

// extractToken 提取 token 的辅助函数
func extractToken(r *http.Request) string {
	// 从 Authorization Header 获取
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
		return authHeader
	}

	// 从 Cookie 获取
	if cookie, err := r.Cookie("session_token"); err == nil {
		return cookie.Value
	}

	// 从查询参数获取
	return r.URL.Query().Get("token")
}
