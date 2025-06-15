package master

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"slices"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/internal/command"
	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/replication"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage/memory/kvstore"
	filemanager "github.com/codecrafters-io/redis-starter-go/app/internal/storage/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/pkg/errors_r"
)

// 确保MasterServer实现replication.MasterServerInterface接口（编译时检查）
var _ replication.MasterServerInterface = (*MasterServer)(nil)

type MasterServer struct {
	*server.BaseServer // cfg & store & registry
	Replicas           []*replicaInfo

	Mu sync.RWMutex
}

type replicaInfo struct {
	conn net.Conn
	addr string
	mu   sync.Mutex // 保护conn的并发访问
}

func NewMasterServer(cfg *config.ServerConfig) *MasterServer {
	store := kvstore.NewStore()
	ms := &MasterServer{
		BaseServer: server.NewBaseServer(cfg, store),
		Replicas:   make([]*replicaInfo, 0), // 初始化为空
	}
	ms.RegisterCmd()
	return ms
}

func (m *MasterServer) RegisterCmd() {
	m.Registry.Register(command.NewPingCommand())
	m.Registry.Register(command.NewEchoCommand())
	// 注册命令
	m.Registry.Register(command.NewSetCommand(m.Store, m.Cfg.Fn, m))
	m.Registry.Register(command.NewGetCommand(m.Store, m.Cfg.Fn))
	m.Registry.Register(command.NewConfigCommand(m.Cfg))
	m.Registry.Register(command.NewKeysCommand(m.Store, m.Cfg.Fn))
	m.Registry.Register(command.NewInfoCommand(m.Cfg))
	m.Registry.Register(command.NewReplconfCommand(m.Cfg))
	m.Registry.Register(command.NewPsyncCommand(m))
}

func (m *MasterServer) Start() error {
	l, err := net.Listen("tcp", ":"+m.Cfg.Port)
	if err != nil {
		log.Printf("Failed to bind to port %s : %s", m.Cfg.Port, err)
		return err
	}
	defer l.Close()

	log.Printf("Master server started on port %s", m.Cfg.Port)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		go m.HandleConnection(conn)
	}
}

func (m *MasterServer) HandleConnection(conn net.Conn) {
	defer func() {
		m.RemoveReplica(conn)
		conn.Close()
	}()
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
		if err := m.ProcessCommand(rw, cmd, args); err != nil {
			log.Printf("Command error: %v, cmd : %s", err, cmd)
			return
		}
	}
}

func (m *MasterServer) ProcessCommand(rw protocol.ResponseWriter, cmd string, args []string) error {
	// 查找命令处理器
	handler, ok := m.Registry.GetHandler(cmd)
	if !ok {
		log.Printf("ERR unknown command '%s'", cmd)
		return errors_r.ErrInvalidRequest
	}

	// 执行命令
	ctx := context.Background()
	return handler.Execute(ctx, rw, args)
}

func (m *MasterServer) AddReplica(conn net.Conn) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	info := &replicaInfo{conn: conn, addr: conn.RemoteAddr().String()}
	m.Replicas = append(m.Replicas, info)

	log.Printf("New replica connected: %s (Total: %d)", info.addr, len(m.Replicas))

	go m.syncToReplica(info)
}

// 删除副本   unused
func (m *MasterServer) RemoveReplica(conn net.Conn) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	if m.Replicas == nil {
		return
	}
	for i, r := range m.Replicas {
		if r.conn == conn {
			// 从切片中移除
			m.Replicas = slices.Delete(m.Replicas, i, i+1)
			log.Printf("Replica disconnected: %s (Remaining: %d)", r.addr, len(m.Replicas))
			return
		}
	}
}

func (m *MasterServer) syncToReplica(info *replicaInfo) {
	info.mu.Lock()
	defer info.mu.Unlock()
	// 实现数据同步逻辑
	if err := m.SendRDBFile(info.conn); err != nil {
		log.Printf("Sync to replica Error: %s", err)
		return
	}
	log.Printf("Sync to replica Success: %s", info.addr)
}

func (m *MasterServer) SendRDBFile(conn net.Conn) error {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	return IoCopyEmpty(m.Cfg.Fn, conn)
}

func IoCopyEmpty(fn string, conn net.Conn) error {
	// REWrite by minimal RDB empty file instead ofnerated from SaveToRDB func. ge
	// minimalRDB :=
	err := filemanager.SaveToRDB(fn, nil)
	if err != nil {
		log.Printf("Save Empty file Error: %s", err)
		return err
	}
	file, err := os.Open(fn)
	if err != nil {
		log.Printf("Open File Error: %s", err)
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	fileSize := fileInfo.Size()
	s := fmt.Sprintf("$%d\r\n", fileSize)
	conn.Write([]byte(s))
	_, err = io.Copy(conn, file) // avoid MemCopy
	if err != nil {
		log.Printf("ioCopy Error: %s", err)
		return err
	}
	return nil
}

func (m *MasterServer) PropagateToReplicas(args []string) error {
	var wg sync.WaitGroup
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	if len(m.Replicas) != 0 {
		log.Printf("Propagating command to %d replicas", len(m.Replicas))
		for idx, c := range m.Replicas {
			if c == nil {
				log.Printf("Warning: nil replica found at index %d", idx)
				continue
			}
			// 可能会有竞态  循环与goroutine异步 可能在goroutine启动时已经进行了几次循环，所以闭包中取值显式赋值
			res := args
			wg.Add(1)
			go func(c *replicaInfo) {
				defer wg.Done()
				if c == nil || c.conn == nil {
					log.Println("Warning: replica or replica.conn is nil")
					return
				}
				if err := c.Write(res); err != nil {
					log.Printf("Propogated Error %s :%s", res[0], err)
					return
				}
				log.Printf("Propogated Success %s, Target Addr:%s", res[0], c.addr)
			}(c)
		}
	} else {
		log.Printf("replConnPool is nil")
	}
	wg.Wait()
	return nil
}

func (r *replicaInfo) Write(args []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.conn == nil {
		return fmt.Errorf("connection is closed")
	}
	_, err := r.conn.Write(protocol.ArrayFmt(args))
	return err
}

func (m *MasterServer) GetPoolLen() int {
	return len(m.Replicas)
}
