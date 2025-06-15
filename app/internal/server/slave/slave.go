package slave

import (
	"context"
	"io"
	"log"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/command"
	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage/memory/kvstore"
	"github.com/codecrafters-io/redis-starter-go/app/pkg/errors_r"
)

type SlaveServer struct {
	*server.BaseServer // cfg & store & registry
	masterConn         net.Conn
}

func NewSlaveServer(cfg *config.ServerConfig) *SlaveServer {
	store := kvstore.NewStore()

	ss := &SlaveServer{
		BaseServer: server.NewBaseServer(cfg, store),
		masterConn: nil,
	}
	ss.RegisterCmd()
	return ss
}

func (s *SlaveServer) RegisterCmd() {
	// 注册命令
	s.Registry.Register(command.NewSetCommand(s.Store, s.Cfg.Fn, nil))
	s.Registry.Register(command.NewGetCommand(s.Store, s.Cfg.Fn))
	s.Registry.Register(command.NewConfigCommand(s.Cfg))
	s.Registry.Register(command.NewKeysCommand(s.Store, s.Cfg.Fn))
	s.Registry.Register(command.NewInfoCommand(s.Cfg))
}

func (s *SlaveServer) Start() error {
	// 连接到主节点
	replConn, err := net.Dial("tcp", s.Cfg.ReplicaOf.MasterHost+":"+s.Cfg.ReplicaOf.MasterPort)
	if err != nil {
		log.Printf("Failed to dial port %s : %s", s.Cfg.Port, err)
		return err
	}
	s.masterConn = replConn
	// 与主节点握手 建立连接
	s.HandShake(replConn)
	log.Printf("握手完成，接受空文件并开始处理主节点消息")
	// 监听并处理主节点输入
	go s.HandleConnection(replConn)

	// 启动从节点服务器监听
	ln, err := net.Listen("tcp", ":"+s.Cfg.Port)
	if err != nil {
		log.Printf("Failed to bind to port %s : %s", s.Cfg.Port, err)
		return err
	}
	log.Printf("Replica server started on port %s", s.Cfg.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go s.HandleConnection(conn)
	}
}

// 启动时同步进行握手 连接建立后再启动命令处理协程
func (s *SlaveServer) HandShake(conn net.Conn) {
	// 1.发送 PING 命令
	rw := protocol.NewConnResponseWriter(conn)
	rw.WriteSimpleString("ping")
	// 判断接收是否为 PONG
	isNormal(conn, "pong", "主节点未响应 PING，可能已经关闭连接")
	// 2.发送 REPLCONF listening-port
	rw.WriteArray([]string{"replconf", "listening-port", s.Cfg.Port})
	isNormal(conn, "ok", "主节点未响应 REPLCONF listening-port")
	// 3.发送 REPLCONF capa
	rw.WriteArray([]string{"replconf", "capa", "psync2"})
	isNormal(conn, "ok", "主节点未响应 REPLCONF capa")
	// 4.发送 PSYNC
	rw.WriteArray([]string{"psync", "?", "-1"}) // PSYNC replId replOffset
	isNormal(conn, "", "")
	// 5.接收(Skip)主节点空文件
	s.skipRDBFile(conn)
}

// 跳过主节点空文件
func (s *SlaveServer) skipRDBFile(conn net.Conn) {
	for i := 0; i < 2; i++ {
		_, _, _ = protocol.ParseRequest(conn)
	}
	return
}

// Wait to be done  目前replid 与 replOffset 未使用 全量复制
// func rcvSYNCinfo(resp string) {
// 	parts := strings.Fields(resp)
// 	if len(parts) == 0 {
// 		return fmt.Errorf("空响应")
// 	}

// 	// 检查第一个部分
// 	if parts[0] != "FULLRESYNC" {
// 		return fmt.Errorf("期望FULLRESYNC，但收到: %s", parts[0])
// 	}

// 	// 检查是否有足够的参数
// 	if len(parts) < 3 {
// 		return fmt.Errorf("不完整的FULLRESYNC响应")
// 	}
// 	replID := parts[1]
// 	offset := parts[2]
// }

func isNormal(conn net.Conn, expectedresp string, errInfo string) {
	resp, _, err := protocol.ParseRequest(conn)
	if err != nil {
		log.Printf("Protocol error: %v", err)
		return
	}
	if strings.HasPrefix(resp, "FULLRESYNC") {
		// rcvSYNCinfo(resp)
		return
	}
	if strings.EqualFold(resp, expectedresp) == false {
		log.Printf("%s", errInfo)
		return
	}
	log.Printf("接收主节点 %s 命令成功", expectedresp)
}

func (s *SlaveServer) HandleConnection(conn net.Conn) {
	defer conn.Close()

	// 创建响应写入器
	rw := protocol.NewConnResponseWriter(conn)
	for {
		// 解析命令
		cmd, args, err := protocol.ParseRequest(conn)
		if err != nil {
			if err != io.EOF {
				log.Printf("Protocol error: %v", err)
				rw.WriteNull()
			}
			return
		}
		// 处理命令
		if err := s.ProcessCommand(rw, cmd, args); err != nil {
			log.Printf("Command error: %v", err)
			return
		}
	}
}

func (s *SlaveServer) ProcessCommand(rw protocol.ResponseWriter, cmd string, args []string) error {
	// 查找命令处理器
	handler, ok := s.Registry.GetHandler(cmd)
	if !ok {
		log.Printf("Slave Rcv ERR unknown command '%s'", cmd)
		return errors_r.ErrInvalidRequest
	}

	// 执行命令
	ctx := context.Background()
	return handler.Execute(ctx, rw, args)
}

// func (s *SlaveServer) Config() *config.ServerConfig {
// 	return s.Cfg
// }
