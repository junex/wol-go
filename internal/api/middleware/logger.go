package middleware

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"time"
)

// Logger 日志中间件
type Logger struct {
	handler http.Handler
}

// NewLogger 创建日志中间件
func NewLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &Logger{handler: next}
	}
}

// ServeHTTP 实现 http.Handler 接口
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// 创建响应记录器，用于捕获状态码
	recorder := &responseRecorder{ResponseWriter: w, status: 200}

	// 调用下一个处理器
	l.handler.ServeHTTP(recorder, r)

	// 记录请求日志
	duration := time.Since(start)
	log.Printf("%s %s %s %d %v",
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
		recorder.status,
		duration,
	)
}

// responseRecorder 响应记录器，用于捕获状态码
type responseRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader 拦截 WriteHeader 以记录状态码
func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Hijack 实现 http.Hijacker 接口，支持 WebSocket
func (r *responseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hijacker.Hijack()
}
