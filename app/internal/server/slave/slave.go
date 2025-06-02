package slave

import (
	"context"
	"io"
	"log"
	"net"

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
	go s.HandleConnection(replConn)

	// 监听客户端
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
		log.Printf("ERR unknown command '%s'", cmd)
		return errors_r.ErrInvalidRequest
	}

	// 执行命令
	ctx := context.Background()
	return handler.Execute(ctx, rw, args)
}

// func (s *SlaveServer) Config() *config.ServerConfig {
// 	return s.Cfg
// }
