/*
Author: Mathieu Mailhos
Filename: dispatcher.go
Description: Dispatch the different job to the pool of miners
*/

package mining

import (
	"gobtcminer/logger"
	"runtime"
	"time"
)

var chunk_queue_capacity int = 300
var monitor logger.Logger
var psize = poolsize()

//Dispatcher Entity.
//Contains a Pool of chans to send and receive from other miners.
//And a queue of chunks to mine
type Dispatcher struct {
	MiningPool chan chan Chunk
	ChunkQueue chan Chunk
}

//Make new Dispatcher
func NewDispatcher(log logger.Logger) *Dispatcher {
	pool := make(chan chan Chunk, psize)
	chunkqueue := make(chan Chunk, chunk_queue_capacity)
	monitor = log
	monitor.Print("info", "New Dispatcher created")
	return &Dispatcher{MiningPool: pool, ChunkQueue: chunkqueue}
}

//Start the new dispatcher, create the miners, start them and begin dispatching.
func (dispatcher *Dispatcher) Run() {
	for i := 0; i < cap(dispatcher.MiningPool); i++ {
		NewMiner(i, dispatcher.MiningPool).Start()
		monitor.Print("info", "New Miner added to the pool")
	}
	dispatcher.dispatch()
}

//Dispatcher start the counter for monitoring. Waits for chunk and send it to an available miner
func (dispatcher *Dispatcher) dispatch() {
	monitor.Print("info", "Starting time counter")
	monitor.BeginTime = time.Now()
	for {
		select {
		case job := <-dispatcher.ChunkQueue:
			go func(job Chunk) {
				//Get a miner available
				BlockToSend := <-dispatcher.MiningPool
				BlockToSend <- job
				monitor.Print("info", "New Chunk sent to the pool")
			}(job)
		}
	}
}

//Set the number of miners depending on the number of threads of the machine.
//Made to optimize and reduce the overhead on multiplex scheduling
func poolsize() int {
	switch maxprocs := runtime.GOMAXPROCS(0); maxprocs {
	case 1:
		return 1
	case 2:
		return 1
	case 3:
		return 2
	default:
		return maxprocs - 2
	}
}
