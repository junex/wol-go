package handler

import (
	"encoding/json"
	"net/http"

	"github.com/junex/wol-go/internal/service"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// 设置 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400, // 24 小时
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "登录成功",
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 从 Cookie 获取 token
	cookie, err := r.Cookie("session_token")
	if err == nil {
		h.authService.Logout(cookie.Value)
	}

	// 清除 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "登出成功",
	})
}

// Status 检查认证状态
func (h *AuthHandler) Status(w http.ResponseWriter, r *http.Request) {
	// 从 Cookie 获取 token
	cookie, err := r.Cookie("session_token")
	token := ""
	if err == nil {
		token = cookie.Value
	}

	// 验证 token
	session, err := h.authService.ValidateToken(token)
	if err != nil {
		h.writeJSON(w, http.StatusOK, map[string]interface{}{
			"success":      false,
			"authenticated": false,
		})
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":       true,
		"authenticated": true,
		"data": map[string]interface{}{
			"username": session.Username,
		},
	})
}

// writeJSON 写入 JSON 响应
func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError 写入错误响应
func (h *AuthHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"message": message,
		},
	})
}
