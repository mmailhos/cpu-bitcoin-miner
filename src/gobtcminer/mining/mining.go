/*
Author: Mathieu Mailhos
Filename: minin.go
Description: Functions for mining a Block Header
*/

package mining

import (
	"errors"
	"gobtcminer/block"
	"strconv"
)

//Macros
const MAX_NONCE uint32 = 4294967295

type ChannelCheck struct {
	Block block.BlockHeader
	Nonce uint32
	Error error
}

//Mining a blockheader and returning the nonce value if suceeded
func Mining_BlockHeader(id int, difficulty float64, bh block.BlockHeader, check_chan chan ChannelCheck) (uint32, error) {
	target := Gettarget(difficulty, bh.Bits)
	for nonce := uint32(0); nonce < MAX_NONCE; nonce++ {
		bh.Nonce = nonce
		if hash := block.Doublesha256_BlockHeader(bh); hash < target {
			check_chan <- ChannelCheck{Block: bh, Nonce: nonce, Error: nil}
			return nonce, nil
		}
	}
	return 0, errors.New("MAX_NONCE reached")
}

//TODO Calculate the right target depending on the difficulty. This one is totally made up for testing purpose.
func Gettarget(difficulty float64, bits uint32) string {
	const padding int = 17
	var target = ""
	for i := 0; i < padding; i++ {
		target = target + "0"
	}
	target = target + strconv.Itoa(int(uint32(difficulty)*bits))
	for i := len(target); i < 64; i++ {
		target = target + "f"
	}
	return target
}
