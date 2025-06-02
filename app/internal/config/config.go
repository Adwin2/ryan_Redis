package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type ServerConfig struct {
	Dir        string
	Dbfilename string
	Fn         string

	Port      string
	Role      string
	ReplicaOf ReplicaConfig
}

type ReplicaConfig struct {
	MasterHost string
	MasterPort string
}

// flag 解析
func ParseFlags() *ServerConfig {
	var cfg ServerConfig
	var replicaOf string

	// 定义命令行参数并绑定到结构体
	pflag.StringVar(&cfg.Dir, "dir", "/var/lib/rdb", "持久化数据存储目录")
	pflag.StringVar(&cfg.Dbfilename, "dbfilename", "dump.rdb", "数据库文件名")
	cfg.Fn = cfg.Dir + "/" + cfg.Dbfilename
	pflag.StringVar(&replicaOf, "replicaof", "", "配置为该地址的副本: '<MASTER_HOST> <MASTER_PORT>'")
	// 为port支持POSIX风格
	pflag.StringVarP(&cfg.Port, "port", "p", "6379", "绑定端口号")
	// 解析参数
	pflag.Parse()
	if replicaOf != "" {
		parts := strings.Fields(replicaOf)
		if len(parts) != 2 {
			log.Println("Error: --replicaof 格式: <MASTER_HOST> <MASTER_PORT>")
			os.Exit(1)
		}

		cfg.ReplicaOf.MasterHost = parts[0]
		cfg.ReplicaOf.MasterPort = parts[1]
		cfg.Role = "slave"
	} else {
		cfg.Role = "master"
	}
	return &cfg
}
