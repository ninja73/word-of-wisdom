package server

import (
	"encoding/hex"
	"math/rand"
	"time"
)

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)

	return hex.EncodeToString(b)[:length]
}
