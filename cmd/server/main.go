package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/junex/wol-go/internal/api"
	"github.com/junex/wol-go/internal/config"
	"github.com/junex/wol-go/internal/repository"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	fmt.Println("======================================")
	fmt.Println("       WOLGO - Wake On LAN Manager")
	fmt.Println("======================================")
	log.Printf("配置加载成功:")
	log.Printf("  端口: %d", cfg.Port)
	log.Printf("  时区: %s", cfg.TZ)
	log.Printf("  认证: %v", cfg.EnableLogin)
	log.Printf("  数据路径: %s", cfg.DataPath)
	log.Printf("  Cron 路径: %s", cfg.CronPath)
	fmt.Println("======================================")

	// 初始化数据访问层
	compRepo := repository.NewComputerRepository(cfg.DataPath + "/computers.txt")
	cronRepo := repository.NewCronRepository(cfg.CronPath + "/wol")

	// 加载计算机数据
	if _, err := compRepo.LoadAll(); err != nil {
		log.Printf("警告: 加载计算机数据失败: %v", err)
	}

	// 初始化路由器
	router := api.NewRouter(cfg)
	router.SetupRoutes()

	// 设置静态文件处理器
	router.SetStaticFileHandler(api.StaticFileHandler(), api.IndexHTML())

	// 启动 HTTP 服务器
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("HTTP 服务器启动在 http://0.0.0.0%s", addr)
	log.Printf("Web 界面: http://0.0.0.0%s", addr)

	// 存储仓库供后续使用（TODO: 通过依赖注入传递）
	_ = compRepo
	_ = cronRepo

	if err := http.ListenAndServe(addr, router.GetRouter()); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
