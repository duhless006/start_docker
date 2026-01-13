package worker

import (
	"log"
	"net/http"
	"sync"
)

type Porkers struct {
	workers map[int]Worker
	mtx     sync.Mutex
}

func NewWorkers() *Porkers {
	return &Porkers{
		workers: make(map[int]Worker),
	}
}

func (p *Porkers) AddWorker(work Worker) error {
	if _, ok := p.workers[work.Id]; ok {
		return http.ErrBodyNotAllowed
	}
	p.workers[work.Id] = work
	return nil

}

func (p *Porkers) GetWorker() (map[int]Worker, error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	tmp := make(map[int]Worker)

	for k, v := range p.workers {
		tmp[k] = v
	}
	return tmp, nil
}

func (p *Porkers) DeleteWorker(id int) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if _, ok := p.workers[id]; !ok {
		return ErrWorkerNotFound
	}
	log.Printf("Deleting worker with ID %d", id)
	delete(p.workers, id)
	return nil
}
