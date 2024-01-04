package helpers

import (
	"crypto/rand"
	"fmt"
)

func Uuid() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)

	if err != nil {
		panic(err)
	}

	// 设置UUID版本为4（基于随机数生成的UUID）
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	// 设置UUID变体为RFC 4122
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	// 将UUID转换为字符串表示
	uuidString := fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])

	return uuidString
}
