/*
Author: Mathieu Mailhos
Filename: mining.go
Description: Functions for mining a Block Header
*/

package mining

import (
	"gobtcminer/block"
	"strconv"
	"time"
)

//MaxNonce standard value
const MaxNonce uint32 = 4294967295

//HashCountSpan counter big enough to avoid mutex bottleneck
const HashCountSpan uint32 = 200000

//Miner entity defined by an ID. Worker.
type Miner struct {
	ID              int
	MiningPool      chan chan Chunk
	BlockChannelIn  chan Chunk
	BlockChannelOut chan Chunk
	quit            chan bool
}

//NewMiner Creates Miner 'Worker'
func NewMiner(id int, miningpool chan chan Chunk, outchan chan Chunk) Miner {
	monitor.Print("info", "New Miner created.")
	return Miner{
		ID:              id,
		MiningPool:      miningpool,
		BlockChannelIn:  make(chan Chunk),
		BlockChannelOut: outchan,
		quit:            make(chan bool)}
}

//Start mining: receive block channels and execute them
func (mine Miner) Start() {
	go func() {
		for {
			//We register the mine into the mining pool
			mine.MiningPool <- mine.BlockChannelIn
			monitor.Print("info", "Miner "+strconv.Itoa(mine.ID)+" available.")
			select {
			//We then receive a chunk to work on or we quit
			case job := <-mine.BlockChannelIn:
				monitor.Print("info", "Miner "+strconv.Itoa(mine.ID)+" starts mining.")
				success, chunk := mine.mining(job)
				if success {
					//Send Back to dispatcher for validation, to be sent back to Websocket
					chunk.Valid = true
					monitor.Print("debug", "Verified Chunk")
					mine.BlockChannelOut <- chunk
				}
			case <-mine.quit:
				return
			}
		}
	}()
}

//Stop tells the Miner to stop working
func (mine Miner) Stop() {
	go func() {
		monitor.Print("info", "Mine "+strconv.Itoa(mine.ID)+" stopped.")
		mine.quit <- true
	}()
}

//Mining a blockheader and returning the chunk including proper nonce value if suceeded. Splited into two to avoid useless checks and increment if the monitoring is not activated. Mining can not take more time than 1 second as the block header expires due to the epoch time field changing constantly.
func (mine *Miner) mining(chunk Chunk) (bool, Chunk) {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	if !monitor.Activated {
		for nonce := chunk.StartNonce; nonce < chunk.EndNonce; nonce++ {
			select {
			case <-timeout:
				//Timeout
				return false, chunk
			default:
				//Success
				chunk.Block.Nonce = nonce
				if hash := block.Doublesha256BlockHeader(chunk.Block); hash < chunk.Target {
					return true, chunk
				}
			}
		}
	} else {
		for count, nonce := uint32(0), chunk.StartNonce; nonce < chunk.EndNonce; nonce, count = nonce+1, count+1 {
			select {
			case <-timeout:
				//Timeout
				monitor.Print("info", "Timeout, moving to next block. "+strconv.Itoa(int(count))+" operations done on this block by Miner "+strconv.Itoa(mine.ID))
				return false, chunk
			default:
				//Success
				chunk.Block.Nonce = nonce
				hash := block.Doublesha256BlockHeader(chunk.Block)
                                if hash < chunk.Target {
					monitor.IncrementBlockCount()
					monitor.Print("info", "NEW BLOCK FOUND!! Nonce:"+strconv.Itoa(int(nonce))+" Miner:"+strconv.Itoa(mine.ID)+" Hash:"+hash)
					return true, chunk
				}
				monitor.Print("debug", "Nonce:"+strconv.Itoa(int(nonce))+" Miner:"+strconv.Itoa(mine.ID)+" Hash:"+hash)
				if count == HashCountSpan {
					monitor.IncrementHashCount(count)
					count = 0
				}
			}
		}
	}
	//MaxNonce Reached
	return false, chunk
}
