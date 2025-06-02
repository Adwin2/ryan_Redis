// command/set.go
package command

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/replication"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage/memory/kvstore"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/pkg/errors_r"
)

type SetCommand struct {
	store  *kvstore.Store
	fn     string
	master replication.MasterServerInterface
}

func NewSetCommand(store *kvstore.Store, fn string, master replication.MasterServerInterface) *SetCommand {
	return &SetCommand{store: store, fn: fn, master: master}
}

func (c *SetCommand) Name() string {
	return "SET"
}

//	func (c *SetCommand) Arity() int {
//		return -3 // 至少需要3个参数: SET key value [EX seconds|PX milliseconds]
//	}
func (c *SetCommand) Execute(ctx context.Context, rw protocol.ResponseWriter, args []string) error {
	if len(args) < 2 {
		return errors_r.ErrWrongNumberOfArguments
	}

	key, value := args[1], args[2]
	// SET 命令
	log.Printf("Received SET, setting %s to %s", key, value)
	// SET... PX... 命令
	if len(args) > 3 && strings.EqualFold(args[3], "PX") == true {
		// 获取设定的过期时间
		exTime, _ := strconv.Atoi(args[4])
		c.store.SetWithExpire(key, value, time.Duration(exTime)*time.Millisecond)
		_ = rdb.UpdateRDB(c.fn, c.store)

		go rdb.Expiry(exTime, c.store, key, c.fn)
	} else {
		// SET
		c.store.Set(key, value)
		_ = rdb.UpdateRDB(c.fn, c.store)
	}

	// 如果是主服务器，传播命令到副本
	if c.master != nil {
		// 传播命令
		if err := c.master.PropagateToReplicas(args); err != nil {
			log.Printf("Failed to propagate SET command: %v", err)
		}
	}

	// 写入OK
	return rw.WriteSimpleString("OK")
}
