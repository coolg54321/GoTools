package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker interface {
	Task()
}

type WorkerItem struct {
	t time.Duration
}

func (wi *WorkerItem) Task() {
	fmt.Println("Going to sleep: ", wi.t)
	time.Sleep(wi.t)
}

type WorkPool struct {
	workers chan Worker
	wg      sync.WaitGroup
}

func NewWorkPool(n int) *WorkPool {
	var w = &WorkPool{
		workers: make(chan Worker),
	}

	w.wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			for p := range w.workers {
				p.Task()
			}
			w.wg.Done()
		}()
	}

	return w
}

func (w *WorkPool) Run(wt Worker) {
	w.workers <- wt
}

func (w *WorkPool) ShutDown() {
	close(w.workers)
	w.wg.Wait()
}
