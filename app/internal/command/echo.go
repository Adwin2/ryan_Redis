// 无需注册

package command

import "github.com/codecrafters-io/redis-starter-go/app/internal/protocol"

func ECHOExecute(rw protocol.ResponseWriter, args []string) error {
	return rw.WriteSimpleString(args[1])
}
