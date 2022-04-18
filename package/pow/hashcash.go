package pow

import (
	"crypto/sha1"
	"fmt"
	"math/big"
)

type hashCash struct {
	BitStrength int32
	Data        string
	Timestamp   int64
	Counter     int64
	Signature   uint64
}

func NewHashCash(bitStrength int32, data string, timestamp int64, signature uint64) *hashCash {
	hc := &hashCash{
		BitStrength: bitStrength,
		Data:        data,
		Timestamp:   timestamp,
		Signature:   signature,
	}

	return hc
}

func (hc hashCash) String() string {
	return fmt.Sprintf(
		"%d:%d:%s:%d:%d",
		hc.BitStrength,
		hc.Timestamp,
		hc.Data,
		hc.Counter,
		hc.Signature,
	)
}

func (hc *hashCash) Check() bool {
	if hc.ZeroCount() >= hc.BitStrength {
		return true
	}
	return false
}

func (hc *hashCash) FindProof() {
	for {
		if hc.Check() {
			return
		}
		hc.Counter++
	}
}

func (hc *hashCash) ZeroCount() int32 {
	digest := sha1.Sum([]byte(hc.String()))
	digestHex := new(big.Int).SetBytes(digest[:])
	return int32((sha1.Size * 8) - digestHex.BitLen())
}
