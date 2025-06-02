package command

import (
	"context"
	"log"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage/memory/kvstore"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage/rdb"
)

type KeysCommand struct {
	store *kvstore.Store
	fn    string
}

func NewKeysCommand(store *kvstore.Store, fn string) *KeysCommand {
	return &KeysCommand{store: store, fn: fn}
}

func (c *KeysCommand) Name() string {
	return "KEYS"
}

func (c *KeysCommand) Execute(ctx context.Context, rw protocol.ResponseWriter, args []string) error {
	var kvs []string
	// KEYS *
	if strings.EqualFold(args[1], "*") == true {
		if c.store.Keys() != nil {
			kvs = c.store.Keys()
		}
		if len(kvs) == 0 {
			kv, err := rdb.GetRDBkeys(c.fn)
			// kv, err := filemanager.TmpParseKV(fn)
			// filemanager.ShowFile(fn)
			if err != nil {
				log.Printf("`KEYS` GetRDBkeys func Wrong: %s", err)
				return err
			}
			// 只输出KEY 偶数位
			for i := 0; i < len(kv); i += 2 {
				kvs = append(kvs, kv[i])
				log.Printf("KEYS WriteArrayFmt: %s", kv[i])
			}
		}
	}
	return rw.WriteArray(kvs)
}
