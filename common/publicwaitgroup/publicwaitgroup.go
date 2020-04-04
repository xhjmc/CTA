package publicwaitgroup

import "sync"

var wg sync.WaitGroup

func Wait() {
	wg.Wait()
}

func Done() {
	wg.Done()
}

func Add(delta int) {
	wg.Add(delta)
}
