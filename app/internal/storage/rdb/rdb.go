package rdb

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"hash/crc64"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/storage/memory/kvstore"
)

// 全局文件锁
var fileLock sync.RWMutex

// header :  "REDIS0011" (fixed)
// metadata: ""   (fixed)  -- > func SaveToRDB()
// DB   KVs (encoded)
// EOF
func SaveToRDB(filename string, store *kvstore.Store) error {
	log.Printf("开始保存文件: %s", filename)
	defer func(start time.Time) {
		log.Printf("文件保存完成: %s (耗时: %v)", filename, time.Since(start))
	}(time.Now())

	fileLock.Lock()
	defer fileLock.Unlock()

	// 确保目录存在
	dir := filepath.Dir(filename)
	ensureDir(dir)
	// if err := os.MkdirAll(dir, 0755); err != nil { // 关键修复点
	// 	return fmt.Errorf("创建目录失败: %v", err)
	// }

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("SaveToRDB Func Openfile Error: %s", err)
		return err
	}
	defer file.Close()

	// 创建CRC64 计算器
	hasher := crc64.New(crc64Table)
	multiWriter := io.MultiWriter(file, hasher) // 同时写到多个Writer 这里用于同时写入文件和计算校验值

	// 文件头
	header := []byte("REDIS0011") // 版本标识
	if _, err := multiWriter.Write(header); err != nil {
		log.Printf("Header Error: %s", err)
		return err
	}

	// 元数据
	if _, err := multiWriter.Write([]byte{
		0xFA,                                                       // 标记元数据区的开始
		0x09, 0x72, 0x65, 0x64, 0x69, 0x73, 0x2D, 0x76, 0x65, 0x72, // "redis-ver"
		0x06, 0x36, 0x2E, 0x30, 0x2E, 0x31, 0x36, // "6.0.16"
	}); err != nil {
		log.Printf("Metadata Error: %s", err)
		return err
	}

	// 数据库部分
	var kvCnt, pxCnt byte = 0, 0
	if store != nil && store.Data != nil && store.Expires != nil {
		kvCnt = byte(len(store.Data))
		pxCnt = byte(len(store.Expires))
	}
	// start:
	if _, err := multiWriter.Write([]byte{
		0xFE, // 标记数据库区的开始
		0x00, // db index
		0xFB, // 标记哈希表大小信息开始
		kvCnt,
		pxCnt,
	}); err != nil {
		log.Printf("DB Start Error: %s", err)
		return err
	}
	if store != nil && store.Data != nil && store.Expires != nil {
		// 遍历处理键值对 并写入
		for key, value := range store.Data {
			// 检测是否需要编码过期时间戳
			if pxTime, exist := store.Expires[key]; exist {
				// uint64可以安全存储13位数
				// 过期时间设置
				pxTimeInt := pxTime.Unix()
				// 持久化到期时间
				timeStamp := uint64(pxTimeInt) + uint64(time.Now().UnixMilli())
				// pxTime int
				multiWriter.Write([]byte{0xFC}) // 标志过期时间戳
				if err := binary.Write(multiWriter, binary.LittleEndian, timeStamp); err != nil {
					log.Printf("Write TimeStamp Error : %s", err)
					return err
				}
			}
			// bytes.Join方法 可能会产生内存峰值(取决于KV大小)
			encodedKV := [][]byte{
				[]byte{0x00}, // Value Type is String
				[]byte{byte(len(key))},
				[]byte(key),
				[]byte{byte(len(value))},
				[]byte(value),
			}
			if _, err := multiWriter.Write(bytes.Join(encodedKV, nil)); err != nil {
				log.Printf("Write len Error: %s", err)
				return err
			}
		}

		// //// ---  delete  ---
		// kvMap.Range(func(key, value any) bool {
		// 	// 处理kv
		// 	k, _ := key.(string)
		// 	v, _ := value.(string)
		// 	// 检测是否需要编码过期时间戳
		// 	if pxTime, exist := pxMap.Load(key); exist {
		// 		// uint64可以安全存储13位数
		// 		timeStamp := pxTime.(uint64) + uint64(time.Now().UnixMilli())
		// 		// pxTime int
		// 		multiWriter.Write([]byte{0xFC}) // 标志过期时间戳
		// 		if err := binary.Write(multiWriter, binary.LittleEndian, timeStamp); err != nil {
		// 			log.Printf("Write TimeStamp Error : %s", err)
		// 			return false
		// 		}
		// 	}
		// 	// bytes.Join方法 可能会产生内存峰值(取决于KV大小)
		// 	encodedKV := [][]byte{
		// 		[]byte{0x00}, // Value Type is String
		// 		[]byte{byte(len(key.(string)))},
		// 		[]byte(k),
		// 		[]byte{byte(len(value.(string)))},
		// 		[]byte(v),
		// 	}
		// 	if _, err := multiWriter.Write(bytes.Join(encodedKV, nil)); err != nil {
		// 		log.Printf("Write len Error: %s", err)
		// 		return false
		// 	}
		// 	return true
		// })
	}
	// 写入结束标记
	multiWriter.Write([]byte{0xFF})

	// 写入8字节(uint64) CRC64校验和
	return binary.Write(file, binary.LittleEndian, hasher.Sum64())
}

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func UpdateRDB(filename string, store *kvstore.Store) error {
	if err := SaveToRDB(filename, store); err != nil {
		log.Printf("UpdateRDB Save File Error: %s", err)
		return err
	}
	return nil
}

// 过期时间Ticker  传入毫秒
func Expiry(t int, store *kvstore.Store, s string, filename string) {
	time.Sleep(time.Duration(t) * time.Millisecond)
	delete(store.Data, s)
	delete(store.Expires, s)
	UpdateRDB(filename, store)
}

func GetRDBkeys(filename string) ([]string, error) {
	var res []string
	var kn int

	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Get Keys Error: %s", err)
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	// n, err := io.CopyN(io.Discard, reader, 32) // 废案
	// FE 数据库部分的标识
	// 找FB KV信息标识
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if b == 0xFB {
			KeyNum, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			kn = int(KeyNum)
			_, _ = reader.ReadByte() // 跳过 PXKV数
			break
		}
	}

	for range kn {
		// 检测 Flag 值, 分别进入不同情况的解析函数
		flag, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Reading bytes Error: %s", err)
			return nil, err
		}
		if flag == 0x00 {
			res = append(res, ParseKV(reader)...)
		} else if flag == 0xFC {
			ts, kv := ParsePXKV(reader)
			// 过期kv会返回nil
			if kv != nil {
				log.Printf("key:%s的时间戳为:%v", kv[0], ts)
				res = append(res, kv...)
			}
			// log.Printf("现在UnixMi值为：%v", time.Now().UnixMilli())
		} else if flag == 0xFF {
			break
		} else {
			log.Printf("unknown flag:%x", flag)
		}
	}
	return res, nil
}

// FC 解析时间戳 KV
func ParsePXKV(reader *bufio.Reader) (ts uint64, res []string) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		panic(err)
	}
	ts = binary.LittleEndian.Uint64(buf)
	_, _ = reader.ReadByte() // 跳过0x00标识 暂时只有00这一种情况
	res = ParseKV(reader)

	if ts <= uint64(time.Now().UnixMilli()) {
		return ts, nil
	}
	return
}

// 00 解析普通KV
func ParseKV(reader *bufio.Reader) (res []string) {
	keylen, err := reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return
		}
		panic(err)
	}
	keyBytes := make([]byte, keylen)
	_, err = io.ReadFull(reader, keyBytes)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		panic(err)
	}
	res = append(res, string(keyBytes))

	vallen, err := reader.ReadByte()
	valBytes := make([]byte, vallen)
	_, err = io.ReadFull(reader, valBytes)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		panic(err)
	}

	res = append(res, string(valBytes))
	return
}

// func SizeMap(m *sync.Map) byte {
// 	cnt := 0
// 	m.Range(func(key, value any) bool {
// 		cnt++
// 		return true
// 	})
// 	return byte(cnt)
// }

var (
	// bufferPool = sync.Pool{
	// 	New: func() any {
	// 		return new(bytes.Buffer)
	// 	},
	// }
	crc64Table = crc64.MakeTable(crc64.ECMA)
)
