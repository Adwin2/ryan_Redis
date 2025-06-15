package replication

import (
	"net"
)

type MasterServerInterface interface {
	AddReplica(conn net.Conn)
	RemoveReplica(conn net.Conn)
	PropagateToReplicas(args []string) error
	SendRDBFile(conn net.Conn) error
	GetPoolLen() int
}
