package user

import (
	"crypto/md5"
	"encoding/hex"
)

func Hashing(str string) string {
	hash := md5.Sum([]byte(str))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}
