package main

import (
	"log"

	svr "github.com/codecrafters-io/redis-starter-go/app/cmd/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
)

func main() {
	conf := config.ParseFlags()
	server := svr.NewServer(conf)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
	select {} // 阻塞主线程
}
