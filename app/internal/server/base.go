package server

import (
	"github.com/codecrafters-io/redis-starter-go/app/internal/command"
	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal/storage/memory/kvstore"
)

type BaseServer struct {
	Cfg      *config.ServerConfig
	Store    *kvstore.Store
	Registry *command.Registry
}

func NewBaseServer(cfg *config.ServerConfig, store *kvstore.Store) *BaseServer {
	return &BaseServer{
		Cfg:      cfg,
		Store:    store,
		Registry: command.NewRegistry(),
	}
}
