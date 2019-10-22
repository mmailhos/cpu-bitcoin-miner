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

//Validate the syntax of each field. Difficulty is not checked since we might need to check older block. Nonce either since it starts at 0.
func Validate(block Header) bool {
	validVersion := false
	versionList := [3]byte{1, 2, 3}
	for _, version := range versionList {
		if block.Version == version {
			validVersion = true
			break
		}
	}
	if validVersion == false {
		return false
	}

	if len(block.HashPrevBlock) > 32 || len(block.HashMerkRoot) > 32 {
		return false
	}

	//1230940800 is 3th Jan 2009 - First Version of Bitcoin
	if block.Time < bitcoinCreationDate || block.Time > uint32(time.Now().Unix()) {
		return false
	}
	return true
}

//MakeSemiRandomBlockHeader makes a semi-random block header. Uses pre-defined time and version. Faster to generate than fully random blockheader.
func MakeSemiRandomBlockHeader(version byte, currentTime uint32) Header {
	hashprevblock := randStringBytes(64)
	hashmerkroot := randStringBytes(64)
	rand.Seed(int64(currentTime))
	nonce := rand.Uint32()
	bits := rand.Uint32()
	return Header{Version: version, HashPrevBlock: hashprevblock, HashMerkRoot: hashmerkroot, Bits: bits, Time: currentTime, Nonce: nonce}
}

//Return the hex string of a given block header.
func hexHeader(bh Header) string {
	hexVersion := strconv.FormatInt(int64(bh.Version), 16) //Little Endian format already, We keep it that way.
	switch length := len(hexVersion); length {
	case 1:
		hexVersion = "0" + hexVersion + "000000"
	case 2:
		hexVersion = hexVersion + "000000"
	}
	hexTime := strconv.FormatInt(int64(bh.Time), 16)
	hexNonce := strconv.FormatInt(int64(bh.Nonce), 16)
	hexBits := strconv.FormatInt(int64(bh.Bits), 16)
	return hexVersion + bh.HashPrevBlock + bh.HashMerkRoot + hexTime + hexBits + hexNonce
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
	data := []byte(hexHeader(bh))
	hash := hash256(hash256(data))
	return hex.EncodeToString(hash)
}
