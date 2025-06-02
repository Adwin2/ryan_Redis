// command/handler.go
package command

import (
	"context"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
)

// Handler 命令处理器接口
type Handler interface {
	Execute(ctx context.Context, rw protocol.ResponseWriter, args []string) error
	Name() string
	// Arity() int // 参数数量，-1 表示可变参数
}

// // Command 命令注册结构
// type Command struct {
// 	Name    string
// 	Handler Handler
// }

// Registry 命令注册表
type Registry struct {
	commands map[string]Handler
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Handler),
	}
}

func (r *Registry) Register(cmd Handler) {
	r.commands[cmd.Name()] = cmd
}

func (r *Registry) GetHandler(name string) (Handler, bool) {
	cmd, ok := r.commands[strings.ToUpper(name)]
	return cmd, ok
}
