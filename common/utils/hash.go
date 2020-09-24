package utils

import (
	"hash/crc32"
)

func CRC32(str string) int32 {
	r := crc32.ChecksumIEEE([]byte(str))
	return int32(r)
}
