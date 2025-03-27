package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
