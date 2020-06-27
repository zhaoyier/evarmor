package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func calcMD5(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
