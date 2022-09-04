package utils

import (
	"crypto/md5"
	"fmt"
)

func GenerateHash(url string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(url)))
}
