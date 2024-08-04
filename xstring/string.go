package xstring

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unsafe"

	"github.com/dlclark/regexp2"
)

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
	letterBytes   = "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	oneStar   = "*"
	twoStar   = "**"
	threeStar = "***"
	fourStar  = "****"
)

var (
	src  = rand.NewSource(time.Now().UnixNano())
	star = map[int]string{
		16: "********",
		15: "*******",
		14: "******",
		13: "*****",
		12: "****",
		11: "***",
		10: "**",
		9:  "*",
	}

	pwdRegx                 = `^(?=(?:.*\d.*\D|.*\D.*\d|.*[a-zA-Z].*[^\w\s]|.*[^\w\s].*[a-zA-Z]|.*[^\w\s].*\d|.*\d.*[^\w\s])).{8,}$`
	chinaPhoneRegex         = `^(?:(?:\+|00)86)?1(?:3\d{3}|4[5-9]\d{2}|5[0-35-9]\d{2}|6[2567]\d{2}|7[0-8]\d{2}|8[0-9]\d{2}|9[189]\d{2})\d{6}$`
	internationalPhoneRegex = `^(?:\+?)(?:[0-9]\d{1,3})?[ -]?\(?(?:\d{1,4})?\)?[ -]?\d{1,14}$`
)

// 密码检测
func MustCompilePwd(password string) bool {
	reg, _ := regexp2.Compile(pwdRegx, 0)

	m, _ := reg.FindStringMatch(password)
	return m != nil
}

// 手机号检测
func MustCompilePhone(areaCode, password string) bool {
	//国内
	if strings.HasPrefix(areaCode, "86") || strings.HasSuffix(areaCode, "86") {
		reg, _ := regexp2.Compile(chinaPhoneRegex, 0)
		m, _ := reg.FindStringMatch(areaCode + password)
		return m != nil
	}
	//国际
	reg, _ := regexp2.Compile(internationalPhoneRegex, 0)
	m, _ := reg.FindStringMatch(areaCode + password)
	return m != nil
}

// 查找
func Contains(s string, array []string) string {
	for _, v := range array {
		if strings.Contains(s, v+".json") {
			return v
		}
	}
	return ""
}

// 搜索
func SearchStrings(s string, array []string) bool {
	for _, v := range array {
		if s == v {
			return true
		}
	}
	return false
}

// 随机字符串
// n长度
func RandString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// 各种屏蔽加*
func ReplaceStringToStar(input string) string {
	var result string
	if input == "" {
		return threeStar
	}
	if strings.Contains(input, "@") {
		res := strings.Split(input, "@")
		if len(res[0]) < 3 {
			result = threeStar + "@" + res[1]
		} else {
			res2 := Substr(input, 0, 3)
			resString := res2 + threeStar
			result = resString + "@" + res[1]
		}
	} else {
		rgx := regexp.MustCompile(`^1[0-9]\d{9}$`)
		if rgx.MatchString(input) {
			result = Substr(input, 0, 3) + threeStar + Substr(input, 7, 11)
		} else {
			nameRune := []rune(input)
			lens := len(nameRune)
			if lens <= 1 {
				result = threeStar
			} else if lens == 2 {
				result = string(nameRune[:1]) + oneStar
			} else if lens == 3 {
				result = string(nameRune[:1]) + oneStar + string(nameRune[2:])
			} else if lens == 4 {
				result = string(nameRune[:1]) + twoStar + string(nameRune[lens-1:])
			} else if 4 < lens && lens <= 10 {
				result = string(nameRune[:2]) + threeStar + string(nameRune[lens-2:])
			} else {
				result = string(nameRune[:3]) + threeStar + string(nameRune[lens-4:])
			}
		}
	}
	return result
}

// 卡号加 *
func CardNumberToStar(cardNumber string) string {
	l := len(cardNumber)
	if l <= 6 {
		return cardNumber
	}
	return cardNumber[:4] + star[l] + cardNumber[len(cardNumber)-4:]
}

// 截取指定长度
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}

// 解析地址
func ParseAddress(srcAddress string) string {
	startIndex := strings.Index(srcAddress, ":")
	endIndex := strings.Index(srcAddress, "?")
	if startIndex > 0 && endIndex > 0 {
		return srcAddress[startIndex+1 : endIndex]
	}
	return "" 
}

// 给地址加*
func ReplaceAddressToStar(address string) string {
	if len(address) > 16 {
		return address[:8] + "..." + address[len(address)-8:]
	}
	return address
}
