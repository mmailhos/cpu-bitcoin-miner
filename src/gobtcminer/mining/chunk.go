/*
Author: Mathieu Mailhos
Filename: chunk.go
Description: Define the chunk entity to mine
*/

package mining

import "gobtcminer/block"

//Chunk entity defined by a block and a target. Nonce is here for checking if the chunk is valid
type Chunk struct {
	Block  block.BlockHeader
	Target string
	Nonce  uint32
}

// Creating Chunk 'job'
func NewChunk(version byte, epoch uint32, difficulty float64) Chunk {
	new_block := block.MakeSemiRandom_BlockHeader(version, epoch)
	target := Gettarget(difficulty, new_block.Bits)
	return Chunk{
		Block:  new_block,
		Target: target,
		Nonce:  0}
}
