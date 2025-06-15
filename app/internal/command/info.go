package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/rtest"
)

type InfoCommand struct {
	Cfg *config.ServerConfig
}

func NewInfoCommand(cfg *config.ServerConfig) *InfoCommand {
	return &InfoCommand{Cfg: cfg}
}

func (c *InfoCommand) Name() string {
	return "INFO"
}

// BulkStringFmt
func (c *InfoCommand) Execute(ctx context.Context, rw protocol.ResponseWriter, args []string) error {
	// return rw.WriteBulkString(fmt.Sprintf("role:master\r\nmaster_repl_offset:%s\r\nmaster_replid:%s\r\n", rtest.Master_repl_offset, rtest.Master_replid))
	if len(args) == 0 {
		return rw.WriteNull()
	}

	section := strings.ToUpper(args[0])

	switch section {
	case "REPLICATION":
		if c.Cfg.Role == "slave" {
			// Handle replica info
			rw.WriteBulkString(fmt.Sprintf("role:slave\r\n"))
			// Add more replica-specific info
		} else {
			rw.WriteBulkString(fmt.Sprintf("role:master\r\nmaster_repl_offset:%s\r\nmaster_replid:%s\r\n", rtest.Master_repl_offset, rtest.Master_replid))
			// Add more master-specific info
		}
		// 其他 待扩展
	}

	return nil
}
