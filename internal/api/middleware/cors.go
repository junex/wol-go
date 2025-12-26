package middleware

import (
	"net/http"
)

// CORS 跨域中间件
type CORS struct {
	handler http.Handler
}

// NewCORS 创建 CORS 中间件
func NewCORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &CORS{handler: next}
	}
}

// ServeHTTP 实现 http.Handler 接口
func (c *CORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 设置 CORS 头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// 处理 OPTIONS 预检请求
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	c.handler.ServeHTTP(w, r)
}
