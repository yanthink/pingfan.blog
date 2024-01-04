package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)

	if err != nil {
		return nil, err
	}

	return key, nil
}

func main() {
	key, err := GenerateRandomKey()
	if err != nil {
		// 处理错误
	}

	fmt.Println(hex.EncodeToString(key)) // 将 key 转换为十六进制字符串输出
}
