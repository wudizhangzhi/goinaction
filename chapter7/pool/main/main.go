package main

import (
	"log"
	"pool"
	"sync/atomic"
	"sync"
	"math/rand"
	"time"
	"io"
)

const (
	maxGoruntines   = 25
	pooledResources = 2
)

type dbConnection struct {
	ID int32
}

func (dbConn dbConnection) Close() error {
	log.Println("Close: Connection: ", dbConn.ID)
	return nil
}

var idCounter int32

func createConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	log.Println("Create: New Connection ", id)
	return &dbConnection{id}, nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(maxGoruntines)

	p, err := pool.New(createConnection, pooledResources)
	if err != nil {
		log.Println(err)
	}

	for query := 0; query < maxGoruntines; query++ {
		go func(q int) {
			performQueries(q, p)
			wg.Done()

		}(query)
	}

	wg.Wait()

	log.Println("Shutdown Program.")
	p.Close()

}

func performQueries(query int, p *pool.Pool) {
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}

	defer p.Release(conn)

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("Query: QID[%d] CID[%d]\n", query, conn.(*dbConnection).ID)

}
