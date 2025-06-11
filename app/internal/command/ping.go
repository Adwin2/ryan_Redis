package command

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
)

type PingCommand struct {
}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

func (c *PingCommand) Name() string {
	return "PING"
}

func (c *PingCommand) Execute(ctx context.Context, rw protocol.ResponseWriter, args []string) error {
	return rw.WriteSimpleString("PONG")
}
