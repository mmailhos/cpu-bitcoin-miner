/*
Author: Mathieu Mailhos
Filename: chunk.go
Description: Define the chunk entity to mine
*/

package mining

import "gobtcminer/block"

//Chunk entity defined by a block and a target. Nonce is here for checking if the chunk is valid
type Chunk struct {
	Block      block.Header
	Target     string
	StartNonce uint32
	EndNonce   uint32
	Valid      bool
}

//NewChunkList We are splitting a BlockHeader to spread the mining between the different goroutines
func NewChunkList(version byte, epoch uint32, difficulty float64) []Chunk {
	newBlock := block.MakeSemiRandomBlockHeader(version, epoch)
	target := Gettarget(difficulty, newBlock.Bits)
	chunklist := make([]Chunk, Psize)
	for i := 0; i < Psize; i++ {
		// Convert values to match MaxNonce type
		poolsize := uint32(Psize)
		iterator := uint32(i)
		startvalue := iterator * (MaxNonce / poolsize)
		if i > 0 {
			startvalue++
		}
		endvalue := (iterator + 1) * (MaxNonce / poolsize)
		chunklist[i] = newChunk(newBlock, target, startvalue, endvalue)
	}
	return chunklist
}

// Returns a new object Chunk
func newChunk(blockheader block.Header, target string, startnonce, endnonce uint32) Chunk {
	return Chunk{
		Block:      blockheader,
		Target:     target,
		StartNonce: startnonce,
		EndNonce:   endnonce,
		Valid:      false,
	}
}
