/*
Author: Mathieu Mailhos
Filename: init.go
Description: Determin the performances of the machin. Gives the machin 5 seconds to mine as many blocks as possible on a given random chunk. This will allow a more efficient split for the block header between the different miners, using this pre-defined value

DEPRECATED - Use of a Timeout per miner instead of this approximation
*/

package initialization

import (
	"fmt"
	"gobtcminer/block"
	"gobtcminer/mining"
	"log"
	"strconv"
	"sync"
	"time"
)

var globalCounter counter

type counter struct {
	HashCount uint32
	mux       sync.Mutex
}

//HashCountSpan Counter big enough to avoid mutex bottleneck but small enough to be as precise as possible
const HashCountSpan = 5000

//Determine the number of double sha256 hash the machin is able to run per second. This function will allow to better split the Chunks with a pre-defined value.
func init() uint32 {
	log.Println("Starting initialization. Determining machin performances... (5s)")
	bh := block.MakeSemiRandomBlockHeader(0, 0)
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(5 * time.Second)
		timeout <- true
	}()
	for i := 0; i < mining.Psize; i++ {
		go noendMining(bh, timeout)
	}
	for {
		select {
		case <-timeout:
			for i := 0; i < mining.Psize; i++ {
				timeout <- true
			}
			fmt.Println("")
			total := globalCounter.HashCount / 5 / uint32(mining.Psize)
			log.Println("Done initializing: " + strconv.Itoa(int(total)) + " operations per seconds per thread")
			return total
		case <-time.After(time.Second * 1):
			fmt.Print(". ")
		}
	}
}

//Miner mining a fake random Chunk while not receiving a timeout
func noendMining(bh block.Header, timeout chan bool) {
	for count, nonce := uint32(0), uint32(0); true; nonce, count = nonce+1, count+1 {
		select {
		case <-timeout:
			return
		default:
			bh.Nonce = nonce
			block.Doublesha256BlockHeader(bh)
			if count == HASHCOUNT_SPAN {
				globalCounter.HashCount += HashCountSpan
				count = 0
			}
		}
	}
}
