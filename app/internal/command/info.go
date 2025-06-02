package command

import (
	"context"
	"fmt"

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
	return rw.WriteBulkString(fmt.Sprintf("role:master\r\nmaster_repl_offset:%s\r\nmaster_replid:%s\r\n", rtest.Master_repl_offset, rtest.Master_replid))
}
