// package socket
//
// //import (
// //	socketio "github.com/googollee/go-socket.io"
// //	"log"
// //)
// //
// //var Server *socketio.Server
// //
// //func InitSocket() *socketio.Server {
// //	server := socketio.NewServer(nil)
// //
// //	server.OnConnect("/", func(s socketio.Conn) error {
// //		log.Println("✅ 客户端连接成功:", s.ID())
// //		return nil
// //	})
// //
// //	server.OnEvent("/", "ping", func(s socketio.Conn, msg string) {
// //		log.Println("收到 ping:", msg)
// //		s.Emit("pong", "服务器已收到")
// //	})
// //
// //	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
// //		log.Println("❌ 客户端断开连接:", s.ID(), reason)
// //	})
// //
// //	Server = server
// //	return server
// //}
//
//	func NotifyAllClients(event string, data interface{}) {
//		//if Server != nil {
//		//	Server.BroadcastToNamespace("/", event, data)
//		//}
//	}
package socket

import (
	socketio "github.com/googollee/go-socket.io"
	"log"
)

var Server *socketio.Server

func InitSocket() *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("✅ 客户端连接成功:", s.ID())
		return nil
	})

	server.OnEvent("/", "ping", func(s socketio.Conn, msg string) {
		log.Println("收到 ping:", msg)
		s.Emit("pong", "服务器已收到")
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("❌ 客户端断开连接:", s.ID(), reason)
	})

	Server = server
	return server
}

func NotifyAllClients(event string, data interface{}) {
	if Server != nil {
		Server.BroadcastToNamespace("/", event, data)
	}
}
