package etag

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/ascii85"
	"fmt"
	"strings"
)

var etagCryptKey32 = []byte("cats_dogs hogetto etag suru desu"[0:32])

func a85fix(r rune) rune {
	switch r {
	case '"':
		return 'v'
	case '\\':
		return 'w'
	case '-':
		return 'x'
	default:
	}
	return r
}
func to_string(src []byte) string {
	dst := make([]byte, ascii85.MaxEncodedLen(len(src)))
	l := ascii85.Encode(dst, src)
	return string(bytes.Map(a85fix, dst[:l]))
}

func Crypt(iv, src []byte) []byte {
	block, err := aes.NewCipher(etagCryptKey32)
	if err != nil {
		panic(fmt.Errorf("fail new cipher: %s", err))
	}
	if len(src) < aes.BlockSize {
		tmp := make([]byte, aes.BlockSize)
		copy(tmp[0:len(src)], src)
		src = tmp
	}

	if len(iv) > aes.BlockSize {
		iv = iv[0:aes.BlockSize]
	} else if len(iv) < aes.BlockSize {
		tmp := make([]byte, aes.BlockSize)
		copy(tmp[0:len(iv)], iv)
		iv = tmp
	}

	dst := make([]byte, len(src))
	copy(iv, src[0:aes.BlockSize])

	ctx := cipher.NewOFB(block, iv)
	ctx.XORKeyStream(dst, src)

	return dst
}

func Make(ids ...[]byte) string {
	etag := make([]string, len(ids))
	for i, id := range ids {
		etag[i] = to_string(id)
	}
	return `"` + strings.Join(etag, "-") + `"`
}

func Split(etag_str string) ([]string, bool) {
	is_weak := false

	etag_str = skipLWS(etag_str)
	if strings.HasPrefix(etag_str, "W/") {
		etag_str = strings.TrimPrefix(etag_str, "W/")
		is_weak = true
	}

	return splitEtag(etag_str), is_weak
}

func skipLWS(str string) string {
	f := findNoneLWS(str)
	return str[f:]
}

func splitEtag(str string) []string {
	var f int
	str = skipLWS(str)

	if str[0] != '"' {
		return nil
	}

	el := make([]string, 0, strings.Count(str, ",")+1)

	for {
		f = findEndEtag(str)
		if f < 0 {
			return nil
		}
		el = append(el, str[:f])
		str = str[f:]

		f = findNextEtag(str)
		if f == 0 {
			break
		}
		if f < 0 {
			return nil
		}
		str = str[f:]
	}

	return el
}

func findNoneLWS(str string) int {
	i := 0
loop:
	for i < len(str) {
		switch str[i] {
		case ' ':
		case '\t':
		default:
			break loop
		}
		i++
	}

	return i
}

func findNextEtag(str string) int {
	i := 0
	var f int
	f = findNoneLWS(str)
	str = str[f:]
	i += f
	if str == "" {
		return 0
	}
	if str[0] != ',' {
		return -1
	}

	str = str[1:]
	i++

	f = findNoneLWS(str)
	str = str[f:]
	i += f
	if str == "" {
		return -2
	}
	if str[0] != '"' {
		return -3
	}

	return i
}

func findEndEtag(str string) int {
	if str == "" {
		return -1
	}
	if str[0] != '"' {
		return -2
	}

	i := 1
	for i < len(str) {
		switch str[i] {
		case '"':
			return i + 1
		case '\\':
			i++
		}

		i++
	}

	return -3
}
