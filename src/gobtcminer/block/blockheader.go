/*
Author: Mathieu Mailhos
Filename: blockheader.go
Description: Functions relative to the block header structure.
For difficulty history, see: http://www.coindesk.com/data/bitcoin-mining-difficulty-time/
*/

package block

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
	"time"
)

const bitcoinCreationDate uint32 = 1230940800

//Header structure
type Header struct {
	Version       byte   //Block Version Number. Decimal format. 4 bytes Little Endian format originally.
	HashPrevBlock string //256bits Hash of the previous block header
	HashMerkRoot  string //256bits Hash on all of the transactions in the block
	Time          uint32 //Timestamp - Epoch time
	Bits          uint32 //Current target calculated with difficulty. Temporary string representation for testing.
	Nonce         uint32 //32Bits number - iterator
}

// Validate a block header
// Difficulty is not checked since we might need to check older block. Nonce neither as it starts at 0.
func (block *Header) Validate() error {
	err := errors.New("invalid version")
	if block.Version < 1 || block.Version > 3 {
		return err
	}

	if len(block.HashPrevBlock) > 32 || len(block.HashMerkRoot) > 32 {
		return err
	}

	//1230940800 is 3th Jan 2009 - First Version of Bitcoin
	if block.Time < bitcoinCreationDate || block.Time > uint32(time.Now().Unix()) {
		return err
	}
	return nil
}

//Return the hex string of a given block header.
func (block *Header) hex() string {
	hexVersion := strconv.FormatInt(int64(block.Version), 16) //Little Endian format already, we keep it that way.
	switch length := len(hexVersion); length {
	case 1:
		hexVersion = "0" + hexVersion + "000000"
	case 2:
		hexVersion = hexVersion + "000000"
	}
	hexTime := strconv.FormatInt(int64(block.Time), 16)
	hexNonce := strconv.FormatInt(int64(block.Nonce), 16)
	hexBits := strconv.FormatInt(int64(block.Bits), 16)
	return hexVersion + block.HashPrevBlock + block.HashMerkRoot + hexTime + hexBits + hexNonce
}

//  Common functions
//

//MakeSemiRandomBlockHeader makes a semi-random block header. Uses pre-defined time and version. Faster to generate than fully random blockheader.
func MakeSemiRandomBlockHeader(version byte, currentTime uint32) Header {
	hashprevblock := randStringBytes(64)
	hashmerkroot := randStringBytes(64)
	rand.Seed(int64(currentTime))
	nonce := rand.Uint32()
	bits := rand.Uint32()
	return Header{Version: version, HashPrevBlock: hashprevblock, HashMerkRoot: hashmerkroot, Bits: bits, Time: currentTime, Nonce: nonce}
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

//Return a Sha256 Hash of given data
func hash256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

//Doublesha256BlockHeader returns a string representation of doubled-hashed block header
func Doublesha256BlockHeader(bh Header) string {
	data := []byte(bh.hex())
	hash := hash256(hash256(data))
	return hex.EncodeToString(hash)
}
