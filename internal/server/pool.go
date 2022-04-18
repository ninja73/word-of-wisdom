package server

import (
	"sync"
)

type bufPool struct {
	pool sync.Pool
}

func (p *bufPool) get() []byte {
	msg := p.pool.Get()
	if msg == nil {
		return make([]byte, 4096)
	}
	return msg.([]byte)
}

func (p *bufPool) put(data []byte) {
	p.pool.Put(data)
}
