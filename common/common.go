package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/spf13/cast"
)

func UnsafeMarshal(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func ToNum(v interface{}) int {
	return int(cast.ToFloat64(fmt.Sprint(v)))
}

func Str(v interface{}) string {
	return fmt.Sprint(v)
}

func LazyLoad[T any](getter func() *T) func() *T {
	var ins *T
	return func() *T {
		if ins == nil {
			ins = getter()
		}
		if ins == nil {
			ins = new(T)
		}
		return ins
	}
}

func Snake2Camel(s string, bigCamel bool) string {
	var buf strings.Builder
	var nextUpper bool
	for idx, r := range s {
		if r == '_' {
			nextUpper = true
			continue
		}

		if idx == 0 && bigCamel {
			buf.WriteRune(unicode.ToUpper(r))
			continue
		}
		if nextUpper {
			buf.WriteRune(unicode.ToUpper(r))
			nextUpper = false
		} else {
			buf.WriteRune(unicode.ToLower(r))
		}
	}
	return buf.String()
}

func Camel2Snake(s string) string {
	var newStr []rune
	for index, v := range []rune(s) {
		if unicode.IsUpper(v) && index != 0 {
			newStr = append(newStr, rune('_'))
		}
		newStr = append(newStr, unicode.ToLower(v))
	}

	return string(newStr)
}

func MD5(source string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(source))
	cipherStr := md5Ctx.Sum(nil)
	return strings.ToLower(hex.EncodeToString(cipherStr))
}
