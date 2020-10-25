package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

// GetAPIKeyHash returns the md5 digest of ts + private key + public key
func GetAPIKeyHash(ts, pvkey, pbkey string) string {
	h := md5.New()
	io.WriteString(h, ts+pvkey+pbkey)
	return fmt.Sprintf("%x", string(h.Sum(nil)))
}
