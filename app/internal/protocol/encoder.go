package protocol

import (
	"strconv"
	"strings"
)

// RESP 简单字符串编码
func SimpleStringFmt(str string) []byte {
	// ByteOutput := []byte("+" + str + "\r\n")
	// return ByteOutput
	var builder strings.Builder
	builder.Grow(len(str) + 16)
	builder.WriteString("+")
	builder.WriteString(str)
	builder.WriteString("\r\n")
	return []byte(builder.String())
}

func BulkStringFmt(str string) []byte {
	// ByteOutput := []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(str), str))
	var builder strings.Builder
	builder.Grow(len(str) + 16)
	builder.WriteString("$")
	builder.WriteString(strconv.Itoa(len(str)))
	builder.WriteString("\r\n")
	builder.WriteString(str)
	builder.WriteString("\r\n")
	return []byte(builder.String())
}

// RESP Array编码
func ArrayFmt(str []string) []byte {
	// log.Printf("长度：%d字符串：%s", len(value), value)
	// Output := fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(name), name, len(value), value)
	// Output := fmt.Sprintf("*%d\r\n", len(str))
	// for i := 0; i < len(str); i++ {
	// 	Output1 := fmt.Sprintf("$%d\r\n%s\r\n", len(str[i]), str[i])
	// 	Output += Output1
	// }
	// ByteOutput := []byte(Output)
	// return ByteOutput
	var builder strings.Builder
	builder.Grow(len(str)*len(str[0]) + 16)
	builder.WriteString("*")
	builder.WriteString(strconv.Itoa(len(str)))
	builder.WriteString("\r\n")
	for i := 0; i < len(str); i++ {
		builder.WriteString("$")
		builder.WriteString(strconv.Itoa(len(str[i])))
		builder.WriteString("\r\n")
		builder.WriteString(str[i])
		builder.WriteString("\r\n")
	}
	return []byte(builder.String())
}

func NullFmt() []byte {
	return []byte("$-1\r\n")
}

// // 弃用  已经封装在connResponseWriter中
// func EncodeSimpleString(conn net.Conn, str string) error {
// 	_, err := conn.Write(SimpleStringFmt(str))
// 	return err
// }

// func EncodeBulkString(conn net.Conn, str string) error {
// 	_, err := conn.Write(BulkStringFmt(str))
// 	return err
// }

// func EncodeArray(conn net.Conn, str []string) error {
// 	_, err := conn.Write(ArrayFmt(str))
// 	return err
// }

// func EncodeNull(conn net.Conn) error {
// 	_, err := conn.Write(NullFmt())
// 	return err
// }
