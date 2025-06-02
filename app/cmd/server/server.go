package server

import (
	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server/master"
	"github.com/codecrafters-io/redis-starter-go/app/internal/server/slave"
)

type Server interface {
	Start() error
	RegisterCmd()
	ProcessCommand(rw protocol.ResponseWriter, cmd string, args []string) error
}

func NewServer(cfg *config.ServerConfig) Server {
	if cfg.Role == "slave" {
		return slave.NewSlaveServer(cfg)
	} else {
		return master.NewMasterServer(cfg)
	}
}
