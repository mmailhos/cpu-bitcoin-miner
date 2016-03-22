/*
Author: Mathieu Mailhos
Filename: blockheader.go
Description: Functions relative to the block header structure.
For difficulty history, see: http://www.coindesk.com/data/bitcoin-mining-difficulty-time/
*/

package block

import "time"

//Macros
const BITCOIN_CREATION_DATE uint32 = 1230940800
const INITIAL_DIFFICULTY float64 = 1

type BlockHeader struct {
	Version       uint32  //Block Version Number
	HashPrevBlock string  //256bits Hash of the previous block header
	HashMerkRoot  string  //256bits Hash on all of the transactions in the block
	Time          uint32  //Timestamp - Epoch time
	Bits          float64 //Current target in compact format
	Nonce         uint32  //32Bits number - iterator
}

//Validate the syntax of each field. Difficulty is not checked since we might need to check older block. Nonce either since it starts at 0.
func Validate(block BlockHeader) bool {
	valid_version := false
	version_list := [3]uint32{02000000, 03000000, 04000000}
	for _, version := range version_list {
		if block.Version == version {
			valid_version = true
			break
		}
	}
	if valid_version == false {
		return false
	}

	if len(block.HashPrevBlock) > 32 || len(block.HashMerkRoot) > 32 {
		return false
	}

	//1230940800 is 3th Jan 2009 - First Version of Bitcoin
	if block.Time < BITCOIN_CREATION_DATE || int64(block.Time) > time.Now().Unix() {
		return false
	}
	if block.Bits < INITIAL_DIFFICULTY {
		return false
	}
	return true
}
