// 主从节点握手
package command

import (
	"context"
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/replication"
	"github.com/codecrafters-io/redis-starter-go/app/internal/rtest"
)

type PsyncCommand struct {
	ms replication.MasterServerInterface
}

func NewPsyncCommand(ms replication.MasterServerInterface) *PsyncCommand {
	return &PsyncCommand{ms: ms}
}

func (c *PsyncCommand) Name() string {
	return "PSYNC"
}

func (c *PsyncCommand) Execute(ctx context.Context, rw protocol.ResponseWriter, args []string) error {
	// return rw.WriteSimpleString("FULLRESYNC" + " " + rtest.FirReplId + " " + rtest.FirReplOffset)
	if len(args) < 2 {
		log.Printf("ERR wrong number of arguments for 'PSYNC' command")
		return nil
	}
	// 发送 FULLRESYNC 响应  格式: +FULLRESYNC <replid> <offset>
	if err := rw.WriteSimpleString("FULLRESYNC" + " " + rtest.FirReplId + " " + rtest.FirReplOffset); err != nil {
		return err
	}
	// 获取连接
	conn := rw.Conn()
	// 添加到副本列表 并发送空RDB文件
	c.ms.AddReplica(conn)
	return nil
}
