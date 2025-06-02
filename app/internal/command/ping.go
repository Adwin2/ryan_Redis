// PING命令默认支持 直接返回PONG字段
package command

import "github.com/codecrafters-io/redis-starter-go/app/internal/protocol"

func PINGExecute(rw protocol.ResponseWriter, args []string) error {
	return rw.WriteSimpleString("PONG")
}
