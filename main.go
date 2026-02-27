package main

import (
	"api/database"
	"api/middlewares"
	"context"
	"embed"
	"errors"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	gin.SetMode(gin.ReleaseMode)
	// 初始化GIN
	r := gin.Default()
	// 设置中间件
	r.Use(middlewares.Cors())

	staticFS, _ := fs.Sub(staticFiles, "static")
	r.StaticFS("/assets", http.FS(mustSub(staticFS, "assets")))
	r.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(staticFS))
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/assets/") {
			c.Status(http.StatusNotFound)
			return
		}
		indexData, err := fs.ReadFile(staticFS, "index.html")
		if err != nil {
			c.String(http.StatusNotFound, "index.html not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexData)
	})

	// 初始化数据库
	err := database.InitDatabase()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 注册路由
	RegisterRoutes(r)

	// 启动Http服务
	srv := &http.Server{
		Addr:    ":30001",
		Handler: r,
	}

	// 启动服务（非阻塞）
	go func() {
		log.Printf("后端服务已启动，监听端口: 30001")
		log.Printf("访问地址：http://127.0.0.1:30001/")
		cmd := exec.Command(`cmd`, `/c`, `start`, `http://127.0.0.1:30001/`)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		_ = cmd.Start()

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 退出处理
	quit := make(chan os.Signal, 1)
	// 监听中断信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	// 5秒超时关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}

	log.Println("服务已正常关闭")
}

func mustSub(f fs.FS, path string) fs.FS {
	s, err := fs.Sub(f, path)
	if err != nil {
		panic(err)
	}
	return s
}
