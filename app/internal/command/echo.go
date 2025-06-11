// 无需注册

package command

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
)

type EchoCommand struct {
}

func NewEchoCommand() *EchoCommand {
	return &EchoCommand{}
}

func (c *EchoCommand) Name() string {
	return "ECHO"
}

func (c *EchoCommand) Execute(ctx context.Context, rw protocol.ResponseWriter, args []string) error {
	return rw.WriteSimpleString(args[1])
}
