package assignment01bca

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Transaction  string
	Nonce        int
	PreviousHash string
	CurrentHash  string
}

type Blockchain struct {
	Blocks []*Block
}

func NewBlock(bc *Blockchain, transaction string, nonce int, previousHash string) (*Block, error) {
	if nonce < 0 {
		return nil, fmt.Errorf("Nonce must be a positive integer")
	}

	if len(previousHash) == 0 {
		if len(bc.Blocks) == 0 {
			// Special case for the genesis block (no previous block)
			previousHash = "0" // Set a default value for the genesis block
		} else {
			return nil, fmt.Errorf("PreviousHash must not be empty")
		}
	}

	// making a new block and calculating hash
	block := &Block{
		Transaction:  transaction,
		Nonce:        nonce,
		PreviousHash: previousHash,
	}
	block.CurrentHash = CalculateHash(block)
	return block, nil
}

func DisplayBlocks(bc *Blockchain) {
	// going through blocks and printing them
	for i, block := range bc.Blocks {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("Transaction: %s\n", block.Transaction)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Previous Hash: %s\n", block.PreviousHash)
		fmt.Printf("Current Hash: %s\n\n", block.CurrentHash)
	}
}

func ChangeBlock(block *Block, newTransaction string) {
	// Change the transaction of a block
	block.Transaction = newTransaction
	block.CurrentHash = CalculateHash(block)
}

func VerifyChain(bc *Blockchain) bool {
	// Verifying integrity of the blockchain
	for i := 1; i < len(bc.Blocks); i++ {
		prevBlock := bc.Blocks[i-1]
		currentBlock := bc.Blocks[i]
		if currentBlock.PreviousHash != prevBlock.CurrentHash {
			return false
		}

		// Verify the current hash of the block
		if CalculateHash(currentBlock) != currentBlock.CurrentHash {
			return false
		}
	}
	return true
}

func CalculateHash(block *Block) string {
	// Calculating hash of the block
	data := fmt.Sprintf("%s%d%s", block.Transaction, block.Nonce, block.PreviousHash)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}
