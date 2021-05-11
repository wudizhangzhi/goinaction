package work

import (
	"log"
	"sync"
)

type Worker interface {
	Task()
}

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func New(maxGoruntines int) *Pool {
	p := Pool{work:make(chan Worker)}
	p.wg.Add(maxGoruntines)

	for i := 0; i < maxGoruntines; i++ {
		go func() {
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}
	return &p
}

// run submit work to the pool
func (p *Pool) Run(w Worker) {
	log.Println("Run work", w)
	p.work <- w
}

func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
