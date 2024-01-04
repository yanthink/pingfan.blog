package config

import (
	"encoding/hex"
)

type jwtConf struct {
	Key []byte
}

var Jwt jwtConf

func loadJwtConfig() {
	hexKey := GetString("jwt.key")
	key, _ := hex.DecodeString(hexKey)
	Jwt.Key = key
}
