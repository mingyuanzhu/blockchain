package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
	"time"
)

/**
Index is the position of the data record in the blockchain
Timestamp is automatically determined and is the time the data is written
BPM or beats per minute, is your pulse rate
Hash is a SHA256 identifier representing this data record
PrevHash is the SHA256 identifier of the previous record in the chain
*/
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PreHash   string
}

var Blockchain []Block

func calculateHash(block Block) (string, error) {
	str := strings.Join([]string{string(block.Index), block.Timestamp, string(block.BPM), block.PreHash}, "_")
	sha256 := sha256.New()
	_, err := sha256.Write([]byte(str))
	if err != nil {
		return "", err
	}
	hash := sha256.Sum(nil)
	return hex.EncodeToString(hash), nil
}

func generateNewBlock(preBlock Block, BPM int) (Block, error) {
	block := Block{
		Index:     preBlock.Index + 1,
		Timestamp: time.Now().String(),
		BPM:       BPM,
		PreHash:   preBlock.Hash,
	}
	hash, err := calculateHash(block)
	if err != nil {
		return block, err
	}
	block.Hash = hash
	return block, nil
}

func isBlockValid(preBlock, newBlock Block) bool {

	if preBlock.Index+1 != newBlock.Index {
		log.Printf("index invalid")
		return false
	}

	if preBlock.Hash != newBlock.PreHash {
		log.Printf("pre hash invalid")
		return false
	}

	hash, err := calculateHash(newBlock)

	if err != nil || newBlock.Hash != hash {
		log.Printf("hash invalid")
		log.Printf(hash)
		log.Printf(newBlock.Hash)
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	mutex.Lock()
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
	mutex.Unlock()
}

func appendChain(block Block) []Block {
	return append(Blockchain, block)
}
