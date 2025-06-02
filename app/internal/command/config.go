package command

import (
	"context"
	"reflect"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
)

type ConfigCommand struct {
	Cfg *config.ServerConfig
}

func NewConfigCommand(cfg *config.ServerConfig) *ConfigCommand {
	return &ConfigCommand{Cfg: cfg}
}

func (c *ConfigCommand) Name() string {
	return "CONFIG"
}

// ArrayFmt
func (c *ConfigCommand) Execute(ctx context.Context, rw protocol.ResponseWriter, args []string) error {
	// CONFIG GET 命令
	if strings.EqualFold(args[1], "GET") == true {
		getName := args[2]
		// 反射获取结构体字段 即配置
		val := reflect.ValueOf(c.Cfg).Elem().FieldByName(getName)
		getvalue := val.String()
		return rw.WriteArray([]string{getName, getvalue})
	}
	return nil
}
