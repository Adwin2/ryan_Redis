package master_test

import (
	"net"
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server/master"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server/slave"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

// server/master/master_integration_test.go
func TestMasterSlaveReplication(t *testing.T) {
	// 启动主服务器
	masterCfg := &config.ServerConfig{Port: "6380", Role: "master"}
	master := master.NewMasterServer(masterCfg)
	go master.Start()

	// 启动从服务器
	slaveCfg := &config.ServerConfig{Port: "6381", Role: "slave", ReplicaOf: config.ReplicaConfig{MasterHost: "localhost", MasterPort: "6380"}}
	slave := slave.NewSlaveServer(slaveCfg)
	go slave.Start()

	// 给主服务器发送命令
	conn, err := net.Dial("tcp", "localhost:6380")
	require.NoError(t, err)
	defer conn.Close()

	// 测试命令传播
	conn.Write([]byte("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	assert.Equal(t, "+OK\r\n", string(buf[:n]))

	// 验证从服务器是否同步
	time.Sleep(100 * time.Millisecond) // 等待同步
	slaveConn, err := net.Dial("tcp", "localhost:6381")
	require.NoError(t, err)
	defer slaveConn.Close()

	slaveConn.Write([]byte("*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"))
	n, _ = slaveConn.Read(buf)
	assert.Equal(t, "$3\r\nbar\r\n", string(buf[:n]))
}
