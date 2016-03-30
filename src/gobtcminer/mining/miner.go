/*
Author: Mathieu Mailhos
Filename: mining.go
Description: Functions for mining a Block Header
*/

package mining

import (
	"errors"
	"gobtcminer/block"
)

//Macros
const MAX_NONCE uint32 = 4294967295

//Mining entity defined by an Id. Worker.
type Miner struct {
	Id           int
	MiningPool   chan chan Chunk
	BlockChannel chan Chunk
	quit         chan bool
}

// Creating Miner 'Worker'
func NewMiner(id int, miningpool chan chan Chunk) Miner {
	return Miner{
		Id:           id,
		MiningPool:   miningpool,
		BlockChannel: make(chan Chunk),
		quit:         make(chan bool)}
}

//Start mining: receive block channels and execute them
func (mine Miner) Start() {
	go func() {
		for {
			//We register the mine into the mining pool
			mine.MiningPool <- mine.BlockChannel
			select {
			//We then receive a chunk to work on or we quit
			case job := <-mine.BlockChannel:
				mine.mining(job)
			case <-mine.quit:
				return
			}
		}
	}()
}

//Tells the Miner to stop working
func (mine Miner) Stop() {
	go func() {
		mine.quit <- true
	}()
}

//Mining a blockheader and returning the nonce value if suceeded
func (mine *Miner) mining(chunk Chunk) (uint32, error) {
	for nonce := uint32(0); nonce < MAX_NONCE; nonce++ {
		chunk.Block.Nonce = nonce
		if hash := block.Doublesha256_BlockHeader(chunk.Block); hash < chunk.Target {
			return nonce, nil
		}
	}
	return 0, errors.New("MAX_NONCE reached")
}
