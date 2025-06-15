package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Dir        string `mapstructure:"dir"`
	Dbfilename string `mapstructure:"dbfilename"`
	Fn         string //计算值 不需要标签

	Port      string        `mapstructure:"port"`
	Role      string        `mapstructure:"role"`
	ReplicaOf ReplicaConfig `mapstructure:"replicaof"`
}

type ReplicaConfig struct {
	MasterHost string `mapstructure:"master_host"`
	MasterPort string `mapstructure:"master_port"`
}

// flag 解析
func LoadConfig() (*ServerConfig, error) {
	// var cfg ServerConfig
	// var replicaOf string

	// 设置默认值
	viper.SetDefault("dir", "/var/lib/rdb")
	viper.SetDefault("dbfilename", "dump.rdb")
	viper.SetDefault("port", "6379")
	viper.SetDefault("role", "master")
	viper.SetDefault("replicaof.master_host", "")
	viper.SetDefault("replicaof.master_port", "")

	// 配置文件查找路径
	viper.AddConfigPath(".")                // main.go 目录
	viper.AddConfigPath("../config")        // app/config/
	viper.AddConfigPath("/etc/ryan_redis/") // 系统配置目录
	viper.AddConfigPath("$HOME/.config/ryan_redis/")

	// 配置文件名（不带扩展名）
	viper.SetConfigName("config")

	// 配置文件类型
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("读取配置文件失败: %v", err)
			return nil, err
		}
		log.Printf("未指定配置文件，使用默认配置: %v", err)
	}
	log.Printf("Using config file: %s", viper.ConfigFileUsed())
	// // 读取环境变量
	// viper.AutomaticEnv()
	// viper.SetEnvPrefix("RR") // 环境变量前缀
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// // 定义命令行参数并绑定到结构体
	// pflag.StringVar(&cfg.Dir, "dir", "/var/lib/rdb", "持久化数据存储目录")
	// pflag.StringVar(&cfg.Dbfilename, "dbfilename", "dump.rdb", "数据库文件名")
	// pflag.StringVar(&replicaOf, "replicaof", "", "配置为该地址的副本: '<MASTER_HOST> <MASTER_PORT>'")
	// // 为port支持POSIX风格
	// pflag.StringVarP(&cfg.Port, "port", "p", "6379", "绑定端口号")
	// 合并命令行参数
	pflag.String("dir", "", "持久化数据存储目录")
	pflag.String("dbfilename", "", "数据库文件名")
	pflag.StringP("port", "p", "", "绑定端口号")
	pflag.String("role", "", "角色：master/slave")
	pflag.String("replicaof", "", "配置为该地址的副本: '<MASTER_HOST> <MASTER_PORT>'")
	// 解析参数
	pflag.Parse()

	//绑定命令行参数
	viper.BindPFlags(pflag.CommandLine)

	// if replicaOf != "" {
	// 	parts := strings.Fields(replicaOf)
	// 	if len(parts) != 2 {
	// 		log.Println("Error: --replicaof 格式: <MASTER_HOST> <MASTER_PORT>")
	// 		os.Exit(1)
	// 	}

	// 	cfg.ReplicaOf.MasterHost = parts[0]
	// 	cfg.ReplicaOf.MasterPort = parts[1]
	// 	cfg.Role = "slave"
	// } else {
	// 	cfg.Role = "master"
	// }
	if replicaOf := viper.GetString("replicaof"); replicaOf != "" {
		parts := strings.Fields(replicaOf)
		log.Printf("replicaof: %s", replicaOf)
		if len(parts) == 2 {
			viper.Set("replicaof.master_host", parts[0])
			viper.Set("replicaof.master_port", parts[1])
			viper.Set("role", "slave")
		}
	}
	// 创建配置结构体
	cfg := &ServerConfig{
		Dir:        viper.GetString("dir"),
		Dbfilename: viper.GetString("dbfilename"),
		Port:       viper.GetString("port"),
		Role:       viper.GetString("role"),
		ReplicaOf: ReplicaConfig{
			MasterHost: viper.GetString("replicaof.master_host"),
			MasterPort: viper.GetString("replicaof.master_port"),
		},
	}

	cfg.Fn = filepath.Join(cfg.Dir, cfg.Dbfilename)
	return cfg, nil
}
