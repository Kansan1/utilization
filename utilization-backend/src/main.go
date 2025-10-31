package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"utilization-backend/config"
	"utilization-backend/src/api/router"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	//socket.InitSocket()

	// 启动 Socket.IO 的后台协程
	//go socket.Server.Serve()
	//defer socket.Server.Close()
	//
	//// 注册事件
	//socket.Server.OnConnect("/", func(s socketio.Conn) error {
	//	log.Println("连接成功:", s.ID())
	//	return nil
	//})
	//
	//socket.Server.OnError("/", func(s socketio.Conn, e error) {
	//	log.Println("错误:", e)
	//})
	//
	//socket.Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
	//	log.Println("断开连接:", s.ID(), reason)
	//})
	//
	//// 启动 goroutine 防止阻塞
	//go socket.Server.Serve()
	//defer socket.Server.Close()

	app := gin.Default()
	//
	//// 注册 socket.io 路由
	//app.GET("/socket.io/*any", gin.WrapH(socket.Server))
	//app.POST("/socket.io/*any", gin.WrapH(socket.Server))
	//
	//app.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"http://localhost:3000"},
	//	AllowMethods:     []string{"GET", "POST"},
	//	AllowHeaders:     []string{"Content-Type"},
	//	AllowCredentials: true,
	//}))

	// 注册其他路由
	router.InitRouter(app)

	if err := app.Run(":9020"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
