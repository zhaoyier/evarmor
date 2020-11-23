package utils

import (
	"hash/crc32"
)

func CRC32(str string) int64 {
	// r := crc32.ChecksumIEEE([]byte(str))
	crc32q := crc32.MakeTable(0xD5828281)
	r := crc32.Checksum([]byte(str), crc32q)
	return int64(r)
}
