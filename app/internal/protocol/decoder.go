// protocol/decoder.go
package protocol

import (
	"io"
	"log"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/pkg/errors_r"
)

// ParseRequest 从网络连接解析 Redis 协议请求
func ParseRequest(conn net.Conn) (string, []string, error) {
	// Read per-connection
	buf := make([]byte, 1024)
	length, err := conn.Read(buf)
	if err != nil {
		log.Printf("Buf read Error: %#v\n", err)
		if err == io.EOF {
			log.Printf("[Error] %s", errors_r.ErrSlaveClosedConn)
			return "", nil, err
		}
		return "", nil, err
	}
	rawdata := string(buf[:length])
	log.Printf("rawdata: %s", rawdata)
	// 去除收尾空格及"\r\n"
	rawdata = strings.TrimSpace(rawdata)
	log.Printf("rawdata after trim: %s", rawdata)
	if rawdata == "" {
		return "", nil, errors_r.ErrInvalidRequest
	}
	res := strings.Split(rawdata, "\r\n")
	log.Println("res:", res)
	// 添加对空数据的检查
	if len(res) == 0 {
		return "", nil, errors_r.ErrInvalidRequest
	}
	// 总个数 + {单个元素的长度 + 元素} * 总个数
	if strings.HasPrefix(res[0], "*") {
		count := (len(res) - 1) / 2 // 去掉总个数除以二
		args := make([]string, count)
		for i, idx := 2, 0; i < len(res); i, idx = i+2, idx+1 {
			args[idx] = res[i]
		}
		log.Printf("len(res): %d, count:%d, args:%v", len(res), count, args)
		return args[0], args, nil
	} else if strings.HasPrefix(res[0], "+") {
		// usually "+PING"
		cmd := strings.TrimSpace(strings.TrimPrefix(res[0], "+"))
		return cmd, []string{cmd}, nil
	}
	return "", nil, errors_r.ErrInvalidRequest
}
