package pow

import (
	"crypto/sha1"
	"fmt"
	"math/big"
)

type HashCash struct {
	BitStrength int32
	Data        string
	Timestamp   int64
	Counter     int64
	Signature   uint64
}

func (hc HashCash) String() string {
	return fmt.Sprintf(
		"%d:%d:%s:%d:%d",
		hc.BitStrength,
		hc.Timestamp,
		hc.Data,
		hc.Counter,
		hc.Signature,
	)
}

func (hc *HashCash) Check() bool {
	if hc.ZeroCount() >= hc.BitStrength {
		return true
	}
	return false
}

func (hc *HashCash) FindProof() {
	for {
		if hc.Check() {
			return
		}
		hc.Counter++
	}
}

func (hc *HashCash) ZeroCount() int32 {
	digest := sha1.Sum([]byte(hc.String()))
	digestHex := new(big.Int).SetBytes(digest[:])
	return int32((sha1.Size * 8) - digestHex.BitLen())
}
