package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"

	"encoding/base64"
)

// 指定Aes加密
func AesEncryptWithIv(origData []byte, key, ivKey string) string {
	k := []byte(key)

	block, err := aes.NewCipher(k)
	if err != nil {
		return ""
	}

	blockSize := block.BlockSize()
	origData = pkcsPadding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(ivKey))
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

// 指定Aes解密
func AesDecryptWithIv(cryted, key, ivKey string) string {
	crytedByte, err := base64.StdEncoding.DecodeString(cryted)
	if err != nil {
		return ""
	}
	// crytedByte := []byte(cryted)
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return ""
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(ivKey))
	orig := make([]byte, len(crytedByte))
	if check(orig, crytedByte, block.BlockSize()) {
		blockMode.CryptBlocks(orig, crytedByte)

		orig = unPadding(orig)
		return string(orig)
	}
	return ""
}

// 从字节aes加密
func AesEncryptWithByte(origData []byte, key string) string {
	k := []byte(key)
	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return ""
	}

	blockSize := block.BlockSize()                            // 获取秘钥块的长度
	origData = pkcsPadding(origData, blockSize)               // 补全码
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize]) // 加密模式
	cryted := make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(cryted, origData)                   // 加密
	return base64.StdEncoding.EncodeToString(cryted)
}

// aes加密(转化成base64)
func AesEncrypt(orig, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return ""
	}

	blockSize := block.BlockSize()                            // 获取秘钥块的长度
	origData = pkcsPadding(origData, blockSize)               // 补全码
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize]) // 加密模式
	cryted := make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(cryted, origData)                   // 加密
	return base64.StdEncoding.EncodeToString(cryted)

}

// aes解密
func AesDecrypt(cryted, key string) string {
	// 转成字节数组
	crytedByte, err := base64.StdEncoding.DecodeString(cryted)
	if err != nil {
		return ""
	}
	k := []byte(key)
	block, err := aes.NewCipher(k) // 分组秘钥
	if err != nil {
		return ""
	}
	blockSize := block.BlockSize()                            // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize]) // 加密模式
	orig := make([]byte, len(crytedByte))

	// 创建数组
	if check(orig, crytedByte, blockSize) {
		blockMode.CryptBlocks(orig, crytedByte) // 解密

		orig = unPadding(orig) // 去补全码
		return string(orig)
	}
	return ""
}

// PKCS7/PKCS5填充
func pkcsPadding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	if padding <= 0 {
		return nil
	}
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7解除
func unPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return nil
	}

	unpadding := int(origData[length-1])
	if length < unpadding {
		return nil
	}
	return origData[:length-unpadding]
}

// zero填充
func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	if padding <= 0 {
		return nil
	}
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

// 检查
func check(dst, src []byte, blockSize int) bool {
	if len(src) == 0 || len(dst) < len(src) || len(src)%blockSize != 0 {
		return false
	}
	return true
}
