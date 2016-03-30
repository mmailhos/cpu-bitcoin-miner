/*
Author: Mathieu Mailhos
Filename: dispatcher.go
Description: Dispatch the different job to the pool of miners
*/

package mining

var ChunkQueueCapacity int = 300

//Dispatcher Entity.
//Contains a Pool of chans to send and receive from other miners.
//And a queue of chunks to mine
type Dispatcher struct {
	MiningPool chan chan Chunk
	ChunkQueue chan Chunk
}

//Make new Dispatcher
func NewDispatcher(max_miners int) *Dispatcher {
	pool := make(chan chan Chunk, max_miners)
	chunkqueue := make(chan Chunk, ChunkQueueCapacity)
	return &Dispatcher{MiningPool: pool, ChunkQueue: chunkqueue}
}

func (dispatcher *Dispatcher) Run() {
	for i := 0; i < cap(dispatcher.MiningPool); i++ {
		NewMiner(i, dispatcher.MiningPool).Start()
	}
	dispatcher.dispatch()
}

//Dispatcher waits for chunk and send it to an available miner
func (dispatcher *Dispatcher) dispatch() {
	for {
		select {
		case job := <-dispatcher.ChunkQueue:
			go func(job Chunk) {
				//Get a miner available
				BlockToSend := <-dispatcher.MiningPool
				BlockToSend <- job

			}(job)
		}
	}
}
