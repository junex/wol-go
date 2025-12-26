package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/junex/wol-go/internal/api/handler"
	"github.com/junex/wol-go/internal/api/middleware"
	"github.com/junex/wol-go/internal/config"
	"github.com/junex/wol-go/internal/repository"
	"github.com/junex/wol-go/internal/service"
)

// Router 路由器
type Router struct {
	router     *mux.Router
	config     *config.Config
	wsHandler  *handler.WebSocketHandler // WebSocket 处理器引用
}

// NewRouter 创建路由器
func NewRouter(cfg *config.Config) *Router {
	r := mux.NewRouter()

	// 应用中间件
	r.Use(middleware.NewRecovery())
	r.Use(middleware.NewLogger())
	r.Use(middleware.NewCORS())

	return &Router{router: r, config: cfg}
}

// SetupRoutes 设置所有路由
func (r *Router) SetupRoutes() {
	// 初始化仓库和服务
	compRepo := repository.NewComputerRepository(r.config.DataPath + "/computers.txt")
	compRepo.LoadAll()

	cronRepo := repository.NewCronRepository(r.config.CronPath + "/wol")

	statusService := service.NewStatusService(
		r.config.PingTimeout,
		r.config.ARPTimeout,
		r.config.TCPTimeout,
	)

	wolService := service.NewWOLService(statusService)
	compService := service.NewComputerService(compRepo)
	cronService := service.NewCronService(cronRepo)
	authService := service.NewAuthService(
		r.config.Username,
		r.config.Password,
		r.config.EnableLogin,
	)

	// 创建 WebSocket Hub
	wshub := service.NewWSHub()
	go wshub.Run() // 启动 Hub 事件循环

	// 创建批量操作服务
	batchService := service.NewBatchService(wolService, statusService)

	// 创建处理器
	compHandler := handler.NewComputerHandler(compService, statusService, wolService)
	networkHandler := handler.NewNetworkHandler(r.config.ARPTimeout, r.config.ARPInterface)
	cronHandler := handler.NewCronHandler(cronService)
	authHandler := handler.NewAuthHandler(authService)
	wsHandler := handler.NewWebSocketHandler(wshub)
	batchHandler := handler.NewBatchHandler(compService, batchService)

	// 保存 WebSocket Handler 引用（用于其他组件广播消息）
	r.wsHandler = wsHandler

	// API 路由
	api := r.router.PathPrefix("/api").Subrouter()

	// 认证相关（无需认证）
	api.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/auth/logout", authHandler.Logout).Methods(http.MethodPost)
	api.HandleFunc("/auth/status", authHandler.Status).Methods(http.MethodGet)

	// 健康检查
	api.HandleFunc("/health", r.healthHandler).Methods(http.MethodGet)

	// WebSocket 连接
	api.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// 设备管理
	api.HandleFunc("/computers", compHandler.ListComputers).Methods(http.MethodGet)
	api.HandleFunc("/computers", compHandler.AddComputer).Methods(http.MethodPost)
	api.HandleFunc("/computers/{mac}", compHandler.UpdateComputer).Methods(http.MethodPut)
	api.HandleFunc("/computers/{mac}", compHandler.DeleteComputer).Methods(http.MethodDelete)
	api.HandleFunc("/computers/{mac}/status", compHandler.GetStatus).Methods(http.MethodGet)
	api.HandleFunc("/computers/{mac}/wake", compHandler.WakeComputer).Methods(http.MethodPost)
	api.HandleFunc("/computers/{mac}/sleep", compHandler.SleepComputer).Methods(http.MethodPost)

	// Cron 任务管理
	api.HandleFunc("/computers/{mac}/crons", cronHandler.GetCrons).Methods(http.MethodGet)
	api.HandleFunc("/computers/{mac}/crons/wake", cronHandler.AddWakeCron).Methods(http.MethodPost)
	api.HandleFunc("/computers/{mac}/crons/sleep", cronHandler.AddSleepCron).Methods(http.MethodPost)
	api.HandleFunc("/computers/{mac}/crons/wake", cronHandler.DeleteWakeCron).Methods(http.MethodDelete)
	api.HandleFunc("/computers/{mac}/crons/sleep", cronHandler.DeleteSleepCron).Methods(http.MethodDelete)

	// 批量操作
	batchHandler.RegisterRoutes(api)

	// 网络扫描
	api.HandleFunc("/network/arp-scan", networkHandler.ARPScan).Methods(http.MethodGet)

	// 静态文件服务（使用嵌入的文件系统）
	// 这个需要在 SetupRoutes 之后设置
}

// SetStaticFileHandler 设置静态文件处理器
func (r *Router) SetStaticFileHandler(staticFS http.FileSystem, indexHTML []byte) {
	// 静态文件服务
	r.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticFS)))

	// 首页（返回 index.html）
	r.router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			// 404，返回首页（SPA 模式）
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	}).Methods(http.MethodGet)
}

// healthHandler 健康检查处理器
func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// GetRouter 获取底层的 mux.Router
func (r *Router) GetRouter() *mux.Router {
	return r.router
}
