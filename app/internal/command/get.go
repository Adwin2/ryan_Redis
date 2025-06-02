package command

import (
	"context"
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage/memory/kvstore"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/pkg/errors_r"
)

type GetCommand struct {
	store *kvstore.Store
	fn    string
}

func NewGetCommand(store *kvstore.Store, fn string) *GetCommand {
	return &GetCommand{store: store, fn: fn}
}

func (c *GetCommand) Name() string {
	return "GET"
}

func (c *GetCommand) Execute(ctx context.Context, rw protocol.ResponseWriter, args []string) error {
	// 从rdb文件中查找
	if len(c.store.Data) == 0 {
		kv, err := rdb.GetRDBkeys(c.fn)
		if err != nil {
			log.Printf("LS HandleCmd `GET` GetRDBkeys func Wrong: %s", err)
			return err
		}
		found := false
		// 在RDB文件的kv组合中查找key
		for i := 0; i < len(kv); i += 2 {
			// GET + 查找键值
			if kv[i] == args[1] {
				log.Printf("GET bulkStringFmt: %s", kv[i+1])
				// conn.Write(rfmt.BulkStringFmt(kv[i+1]))
				found = true
				return rw.WriteBulkString(kv[i+1])
			}
		}
		// 没有对应key
		if !found {
			return errors_r.ErrKeyNotFoundInRDB
		}
	} else {
		// 从Map中查找
		log.Printf("Received GET, getting %s", args[1])
		OP, ok := c.store.Get(args[1])
		if !ok {
			log.Printf("ERROR GET: %s", args[1])
			return errors_r.ErrKeyNotFoundInMap
		}
		return rw.WriteBulkString(OP)
	}
	return errors_r.ErrKeyNotFound
}
