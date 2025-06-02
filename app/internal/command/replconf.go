package command

import (
	"context"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
)

type ReplconfCommand struct {
	Cfg *config.ServerConfig
}

func NewReplconfCommand(cfg *config.ServerConfig) *ReplconfCommand {
	return &ReplconfCommand{Cfg: cfg}
}

func (c *ReplconfCommand) Name() string {
	return "REPLCONF"
}

func (c *ReplconfCommand) Execute(ctx context.Context, rw protocol.ResponseWriter, args []string) error {
	// 设定侦听从节点的端口
	if strings.EqualFold(args[1], "listening-port") == true {
		c.Cfg.Port = args[2]
	}
	return rw.WriteSimpleString("OK")
}
