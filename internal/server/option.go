package server

import (
	"time"
)

const (
	defaultBitStrength = 20
	defaultTimeout     = 5 * time.Second
	defaultExpiration  = 5 * time.Minute
	defaultLimit       = 100
)

var (
	defaultSecretKey = randomString(20)
)

type Options struct {
	BitStrength int32
	Timeout     time.Duration
	SecretKey   string
	Expiration  time.Duration
	Limit       int32
}

func setDefaultWorkOptions(o *Options) {
	if o.Timeout == 0 {
		o.Timeout = defaultTimeout
	}

	if o.SecretKey == "" {
		o.SecretKey = defaultSecretKey
	}

	if o.Limit == 0 {
		o.Limit = defaultLimit
	}

	if o.BitStrength == 0 {
		o.BitStrength = defaultBitStrength
	}

	if o.Expiration == 0 {
		o.Expiration = defaultExpiration
	}
}
