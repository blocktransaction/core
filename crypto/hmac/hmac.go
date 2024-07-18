package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// 加密base64处理
func Crypto(key, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(key))

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// 加密 hex处理
func HexCrypto(key, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(key))

	return hex.EncodeToString(hash.Sum(nil))
}
