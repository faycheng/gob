package worker

import (
	"sync"
)

type Pool interface {
	Get() Runner
	Put(w Runner)
}

type pool struct {
	workers   sync.Pool
	initSize  int
	threshold int
}

// TODO: auto scaling runner size
func NewPool() Pool {
	workers := sync.Pool{
		New: func() interface{} {
			return &runner{}
		},
	}
	return &pool{
		workers: workers,
	}
}

func (p *pool) Get() Runner {
	return p.workers.Get().(Runner)
}

func (p *pool) Put(w Runner) {
	p.workers.Put(w)
}
