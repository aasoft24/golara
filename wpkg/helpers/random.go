// pkg/helpers/random.go
package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}
