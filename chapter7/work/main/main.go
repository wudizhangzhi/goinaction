package main

import (
	"log"
	"time"
	"sync"
	"work"
)

var names = []string{
	"steve",
	"bob",
	"mary",
	"therese",
	"jason",
}

type namePrinter struct {
	name string
}


func (np *namePrinter) Task() {
	log.Println(np.name)
	time.Sleep(time.Second)
}


func main() {
	p := work.New(2)

	var wg sync.WaitGroup

	wg.Add(100*len(names))

	for i:=0;i<100;i++ {
		for _, name := range names {
			n := namePrinter{name:name}
			go func() {
				p.Run(&n)
				wg.Done()
			}()
		}
	}

	wg.Wait()

	p.Shutdown()

}
