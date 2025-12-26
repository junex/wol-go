package middleware

import (
	"log"
	"net/http"
)

// Recovery 错误恢复中间件
type Recovery struct {
	handler http.Handler
}

// NewRecovery 创建错误恢复中间件
func NewRecovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &Recovery{handler: next}
	}
}

// ServeHTTP 实现 http.Handler 接口
func (r *Recovery) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("PANIC: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	r.handler.ServeHTTP(w, req)
}
