package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/blocktransaction/core/crypto/aes"
	"github.com/blocktransaction/core/xtime"
	"github.com/golang-jwt/jwt/v5"
)

type myCustomClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

type Jwt struct {
	aesSecret    string
	jwtSecret    string
	jwtExpiresAt time.Duration
	issuer       string
	expiresAt    float64
}

// new jwt
func NewJwt(aesSecret, jwtSecret, issuer string, jwtExpiresAt time.Duration) *Jwt {
	return &Jwt{
		aesSecret:    aesSecret,
		jwtSecret:    jwtSecret,
		jwtExpiresAt: jwtExpiresAt,
		issuer:       issuer,
	}
}

// 生成jwt
func (j *Jwt) GenerateJwt(content string) (string, error) {
	if content == "" {
		return "", errors.New("content is empty")
	}

	claims := myCustomClaims{
		aes.AesEncrypt(content, j.aesSecret),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.jwtExpiresAt)),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.jwtSecret))
}

// 解析jwT
func (j *Jwt) ParseJwt(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("token is empty")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok && token.Valid {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.jwtSecret), nil
	})

	if err != nil || token == nil {
		return "", fmt.Errorf("token invalid: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		j.expiresAt = claims["exp"].(float64)
		return aes.AesDecrypt(claims["userId"].(string), j.aesSecret), nil
	}
	return "", errors.New("token invalid")

}

// 验证jwt是否过期
func (j *Jwt) Valid() bool {
	return xtime.Second() < int64(j.expiresAt)
}
