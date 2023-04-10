package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"strings"
)

func Authorize(akID, signature string) string {
	return strings.Join([]string{akID, signature}, ":")
}

func Sign(info, secret []byte) []byte {
	return HMAC_SHA1(info, secret)
}

func HMAC_SHA1(src, key []byte) []byte {
	m := hmac.New(sha1.New, key)
	m.Write([]byte(src))
	return m.Sum(nil)
}

func Verify(msg, sign, secret []byte) bool {
	hmac_sha1 := HMAC_SHA1(msg, secret)
	equal := hmac.Equal(sign, hmac_sha1)

	return equal
}
