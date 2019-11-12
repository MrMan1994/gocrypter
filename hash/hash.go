package hash

import (
	"crypto/md5"
	"encoding/hex"
)

// Create a new md5 hash from the provided string
func Create(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
