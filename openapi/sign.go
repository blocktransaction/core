package openapi

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
)

// 验证签名是否正确
func VerificationSign(data interface{}, apiSecret, signValue string) (bool, string) {
	urlValeus := json2UrlValues(data)
	query, err := url.QueryUnescape(urlValeus.Encode())
	if err != nil || query == "" {
		return false, ""
	}
	query = query + "&key=" + apiSecret
	signResult := sign(query)
	return signValue == signResult, signResult
}

// 签名
func Sign(data interface{}, apiSecret string) string {
	urlValeus := json2UrlValues(data)
	query, err := url.QueryUnescape(urlValeus.Encode())
	if err != nil || query == "" {
		return ""
	}
	query = query + "&key=" + apiSecret
	return sign(query)
}

// 签名
// 将字符串进行sha256加密，并转换成32位base64值，并转换成大写
func sign(query string) string {
	return strings.ToUpper(xhmac(query))
}

// md5加密
func xmd5(param string) string {
	h := md5.New()
	h.Write([]byte(param))

	b := h.Sum(nil)

	return hex.EncodeToString(b)
}

// sha256加密
func xhmac(param string) string {
	sha := sha256.New()
	sha.Write([]byte(param))
	return hex.EncodeToString(sha.Sum(nil))
}
