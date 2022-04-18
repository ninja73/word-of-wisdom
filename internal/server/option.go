package server

import "time"

type Option func(*server)

func WithTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.timeout = timeout
	}
}

func WithBitStrength(bitStrength int32) Option {
	return func(s *server) {
		s.bitStrength = bitStrength
	}
}

func WithSecretKey(secretKey string) Option {
	return func(s *server) {
		s.secretKey = secretKey
	}
}
