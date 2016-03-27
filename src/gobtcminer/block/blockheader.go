/*
Author: Mathieu Mailhos
Filename: blockheader.go
Description: Functions relative to the block header structure.
For difficulty history, see: http://www.coindesk.com/data/bitcoin-mining-difficulty-time/
*/

package block

import (
	"math/rand"
	"strconv"
	"time"
)

//Macros
const BITCOIN_CREATION_DATE uint32 = 1230940800
const INITIAL_DIFFICULTY float64 = 1

type BlockHeader struct {
	Version       byte    //Block Version Number. Decimal format. 4 bytes Little Endian format originally.
	HashPrevBlock string  //256bits Hash of the previous block header
	HashMerkRoot  string  //256bits Hash on all of the transactions in the block
	Time          uint32  //Timestamp - Epoch time
	Bits          float64 //Current target (difficulty) in compact format
	Nonce         uint32  //32Bits number - iterator
}

//Validate the syntax of each field. Difficulty is not checked since we might need to check older block. Nonce either since it starts at 0.
func Validate(block BlockHeader) bool {
	valid_version := false
	version_list := [3]byte{1, 2, 3}
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
	if block.Time < BITCOIN_CREATION_DATE || block.Time > uint32(time.Now().Unix()) {
		return false
	}
	if block.Bits < INITIAL_DIFFICULTY {
		return false
	}
	return true
}

//Make a semi-random block header. Uses pre-defined difficulty, time and version. Faster to generate than fully random blockheader.
func MakeSemiRandom_BlockHeader(difficulty float64, version byte, time uint32) BlockHeader {
	hashprevblock := randStringBytes(64)
	hashmerkroot := randStringBytes(64)
	nonce := rand.Uint32()
	return BlockHeader{Version: version, HashPrevBlock: hashprevblock, HashMerkRoot: hashmerkroot, Bits: difficulty, Time: time, Nonce: nonce}
}

//Return the hex string of a given block header.
func Hex_BlockHeader(bh BlockHeader) string {
	hex_version := strconv.FormatInt(int64(bh.Version), 16) //Little Endian format already, We keep it that way.
	switch length := len(hex_version); length {
	case 1:
		hex_version = "0" + hex_version + "000000"
	case 2:
		hex_version = hex_version + "000000"
	}
	hex_time := strconv.FormatInt(int64(bh.Time), 16)
	hex_bits := strconv.FormatInt(int64(uint32(bh.Bits)), 16)
	hex_nonce := strconv.FormatInt(int64(bh.Nonce), 16)
	return hex_version + bh.HashPrevBlock + bh.HashMerkRoot + hex_time + hex_bits + hex_nonce
}

//Generate Hex string-representated number of n characters
func randStringBytes(n int) string {
	const letterBytes = "abcdef0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
