package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/blocktransaction/core/crypto/aes"
	"github.com/golang-jwt/jwt/v5"
)

type myCustomClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

type Jwt struct {
	AesSecret    string
	JwtSecret    string
	JwtExpiresAt int64
	Issuer       string
}

// new jwt
func NewJwt(aesSecret, jwtSecret, issuer string, jwtExpiresAt int64) *Jwt {
	return &Jwt{
		AesSecret:    aesSecret,
		JwtSecret:    jwtSecret,
		JwtExpiresAt: jwtExpiresAt,
		Issuer:       issuer,
	}
}

// 生成jwt
func (j *Jwt) GenerateJwt(userId, mobile string) (string, error) {
	if userId == "" {
		return "", errors.New("userid empty")
	}

	claims := myCustomClaims{
		aes.AesEncrypt(userId+"|"+mobile, j.AesSecret),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.JwtExpiresAt) * time.Hour)),
			Issuer:    j.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.JwtSecret))
}

// 解析jwt
func (j *Jwt) ParseJwt(tokenString string) (string, float64, error) {
	if tokenString == "" {
		return "", 0.0, errors.New("token is empty")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok && token.Valid {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.JwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return aes.AesDecrypt(claims["userId"].(string), j.AesSecret), claims["exp"].(float64), nil
	}
	return "", 0.0, err
}
