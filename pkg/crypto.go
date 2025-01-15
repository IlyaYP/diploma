package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Hash(m, k string) string {
	h := hmac.New(sha256.New, []byte(k))
	h.Write([]byte(m))
	dst := h.Sum(nil)

	//log.Printf("%s:%x", m, dst)
	return hex.EncodeToString(dst)
}
