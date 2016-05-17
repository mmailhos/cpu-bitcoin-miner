/*
Author: Mathieu Mailhos
Filename: chunk.go
Description: Define the chunk entity to mine
*/

package mining

import "gobtcminer/block"

//Chunk entity defined by a block and a target. Nonce is here for checking if the chunk is valid
type Chunk struct {
	Block      block.BlockHeader
	Target     string
	StartNonce uint32
	EndNonce   uint32
}

// We are splitting a BlockHeader to spread the mining between the different goroutines
func NewChunkList(version byte, epoch uint32, difficulty float64) []Chunk {
	new_block := block.MakeSemiRandom_BlockHeader(version, epoch)
	target := Gettarget(difficulty, new_block.Bits)
	chunklist := make([]Chunk, Psize)
	for i := 0; i < Psize; i++ {
		// Convert values to match MAX_NONCE type
		poolsize := uint32(Psize)
		iterator := uint32(i)
		startvalue := iterator * (MAX_NONCE / poolsize)
		if i > 0 {
			startvalue++
		}
		endvalue := (iterator + 1) * (MAX_NONCE / poolsize)
		chunklist[i] = newChunk(new_block, target, startvalue, endvalue)
	}
	return chunklist
}

// Returns a new object Chunk
func newChunk(blockheader block.BlockHeader, target string, startnonce, endnonce uint32) Chunk {
	return Chunk{
		Block:      blockheader,
		Target:     target,
		StartNonce: startnonce,
		EndNonce:   endnonce}
}
