package slip

import (
	"crypto/sha256"
	"encoding/base64"
)

func hash256(key string) string {
	bb := sha256.Sum256([]byte(key))
	return base64.RawURLEncoding.EncodeToString(bb[:])
}
