package command

import (
	"context"
	"log"
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
	if len(args) < 2 {
		log.Printf("[REPLCONF] wrong number of arguments for 'REPLCONF' command")
		return rw.WriteNull()
	}

	//  args : [REPLCONF subcommand arg1 arg2 ...]
	switch strings.ToUpper(args[1]) {
	case "LISTENING-PORT":
		// Store replica's listening port
		c.Cfg.Port = args[2]
		return rw.WriteSimpleString("OK")
	case "CAPA":
		// Handle capabilities
		return rw.WriteSimpleString("OK")
	case "GETACK":
		// Handle ACK from replica
		return rw.WriteSimpleString("OK")
	default:
		log.Printf("[REPLCONF] unknown subcommand '%s'", args[1])
		return rw.WriteNull()
	}
}
