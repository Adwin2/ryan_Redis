package main

import (
	"log"

	svr "github.com/codecrafters-io/redis-starter-go/app/cmd/server"
	"github.com/codecrafters-io/redis-starter-go/app/internal/config"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	server := svr.NewServer(conf)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
	select {} // 阻塞主线程
}
