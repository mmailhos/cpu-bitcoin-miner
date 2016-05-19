/*
Author: Mathieu Mailhos
Filename: mining.go
Description: Functions for mining a Block Header
*/

package mining

import (
	"errors"
	"gobtcminer/block"
	"strconv"
	"time"
)

//Macros
const MAX_NONCE uint32 = 4294967295
const HASHCOUNT_SPAN uint32 = 200000 //Counter big enough to avoid mutex bottleneck

//Mining entity defined by an Id. Worker.
type Miner struct {
	Id           int
	MiningPool   chan chan Chunk
	BlockChannel chan Chunk
	quit         chan bool
}

// Creating Miner 'Worker'
func NewMiner(id int, miningpool chan chan Chunk) Miner {
	monitor.Print("info", "New Miner created.")
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
			monitor.Print("info", "Miner "+strconv.Itoa(mine.Id)+" available.")
			select {
			//We then receive a chunk to work on or we quit
			case job := <-mine.BlockChannel:
				monitor.Print("info", "Miner "+strconv.Itoa(mine.Id)+" starts mining.")
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
		monitor.Print("info", "Mine "+strconv.Itoa(mine.Id)+" stopped.")
		mine.quit <- true
	}()
}

//Mining a blockheader and returning the nonce value if suceeded. Splited into two to avoid useless checks and increment if the monitoring is not activated. Mining can not take more time than 1 second as the block header expires due to the epoch time field changing constantly.
func (mine *Miner) mining(chunk Chunk) (uint32, error) {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	if !monitor.Activated {
		for nonce := chunk.StartNonce; nonce < chunk.EndNonce; nonce++ {
			select {
			case <-timeout:
				return 0, nil
			default:
				chunk.Block.Nonce = nonce
				if hash := block.Doublesha256_BlockHeader(chunk.Block); hash < chunk.Target {
					return nonce, nil
				}
			}
		}
	} else {
		for count, nonce := uint32(0), chunk.StartNonce; nonce < chunk.EndNonce; nonce, count = nonce+1, count+1 {
			select {
			case <-timeout:
				monitor.Print("info", "Timeout, moving to next block. "+strconv.Itoa(int(count))+" operations done on this block by Miner "+strconv.Itoa(mine.Id))
				return 0, nil
			default:

				chunk.Block.Nonce = nonce
				if hash := block.Doublesha256_BlockHeader(chunk.Block); hash < chunk.Target {
					monitor.IncrementBlockCount()
					monitor.Print("info", "NEW BLOCK FOUND!! Nonce:"+strconv.Itoa(int(nonce))+" Miner:"+strconv.Itoa(mine.Id)+" Hash:"+hash)
					return nonce, nil
				} else {
					monitor.Print("debug", "Nonce:"+strconv.Itoa(int(nonce))+" Miner:"+strconv.Itoa(mine.Id)+" Hash:"+hash)
				}
				if count == HASHCOUNT_SPAN {
					monitor.IncrementHashCount(count)
					count = 0
				}
			}
		}
	}
	return 0, errors.New("MAX_NONCE reached")
}
