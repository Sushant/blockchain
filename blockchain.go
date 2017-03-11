package main

import (
	"crypto/sha256"
	"fmt"
)

type Sha256Hash [sha256.Size]byte

type BlockChain struct {
	NewestBlockPointer *HashPointer
	Blocks             map[Sha256Hash]Block
	NewestBlockIndex   int
}

type HashPointer struct {
	Prev     *Block
	PrevHash Sha256Hash
}

type Block struct {
	PrevBlock *HashPointer
	Header    Sha256Hash
	Nonce     int
	Payload   string
}

func main() {
	bc := NewBlockChain()
	bc.AddBlock("Genesis Block")
	fmt.Printf("Current Block %x\n", bc.NewestBlockPointer.PrevHash)
	newestBlock := bc.Blocks[bc.NewestBlockPointer.PrevHash]
	fmt.Printf("Block Header %x\n", newestBlock.Header)
	fmt.Printf("Nonce %d\n", newestBlock.Nonce)
}

func NewBlockChain() *BlockChain {
	return &BlockChain{Blocks: make(map[Sha256Hash]Block)}
}

func (bc *BlockChain) AddBlock(payload string) {
	candidateBlock := bc.mine(payload, bc.NewestBlockIndex+1)
	blockHash := BlockHash(candidateBlock)
	bc.Blocks[blockHash] = *candidateBlock
	bc.NewestBlockIndex += 1
	bc.NewestBlockPointer = &HashPointer{Prev: candidateBlock, PrevHash: blockHash}
}

func (bc *BlockChain) mine(payload string, index int) *Block {
	var header Sha256Hash
	nonce := 0
	for {
		trialBlock := Block{Payload: payload, Nonce: nonce, PrevBlock: bc.NewestBlockPointer}
		header = BlockHash(&trialBlock)
		if header[0] == 0 {
			break
		} else {
			nonce += 1
		}
	}
	return &Block{PrevBlock: bc.NewestBlockPointer, Header: header, Nonce: nonce, Payload: payload}
}

func BlockHash(block *Block) Sha256Hash {
	serializedBlock := fmt.Sprintf("%v", block)
	return sha256.Sum256([]byte(serializedBlock))
}
