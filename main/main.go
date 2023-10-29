package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	Transactions        []string
	Nonce               int
	PreviousHash        string
	CurrentHash         string
	MerkleRoot          string
	TransactionDateTime time.Time
}

type Blockchain struct {
	Blocks                  []*Block
	NumTransactionsPerBlock int
	MinBlockHash            string
	MaxBlockHash            string
}

func NewBlock(bc *Blockchain, transactions []string, nonce int, previousHash string, transactionDateTime time.Time) (*Block, error) {
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

	block := &Block{
		Transactions: transactions,
		Nonce:        nonce,
		PreviousHash: previousHash,
	}
	block.MerkleRoot = CalculateMerkleRoot(transactions)
	block.CurrentHash = CalculateHash(block)

	block.TransactionDateTime = transactionDateTime

	return block, nil
}

func DisplayBlocks(bc *Blockchain) {
	for i, block := range bc.Blocks {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("Transactions: %v\n", block.Transactions)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Previous Hash: %s\n", block.PreviousHash)
		fmt.Printf("Merkle Root: %s\n", block.MerkleRoot)
		fmt.Printf("Current Hash: %s\n", block.CurrentHash)
		fmt.Printf("Transaction Date and Time: %s\n\n", block.TransactionDateTime)
	}
}

func ChangeBlock(block *Block, newTransaction string, place int) {
	block.Transactions[place] = newTransaction
	block.CurrentHash = CalculateHash(block)
}

func VerifyChain(bc *Blockchain) bool {
	// Check if the blockchain is empty
	if len(bc.Blocks) == 0 {
		return false
	}

	// Check if the Genesis block has the correct properties
	genesisBlock := bc.Blocks[0]
	if genesisBlock.PreviousHash != "0" || genesisBlock.Nonce != 0 {
		return false
	}

	calculatedMerkleRoot := CalculateMerkleRoot(genesisBlock.Transactions)
	if calculatedMerkleRoot != genesisBlock.MerkleRoot {
		println("Merkle root mismatch")
		return false
	}

	for i := 1; i < len(bc.Blocks); i++ {
		prevBlock := bc.Blocks[i-1]
		currentBlock := bc.Blocks[i]

		// Check if the PreviousHash of each block matches the CurrentHash of the previous block
		if currentBlock.PreviousHash != prevBlock.CurrentHash {
			return false
		}

		// Check if the nonce is a positive integer
		if currentBlock.Nonce < 0 {
			return false
		}

		// Check if the Merkle root matches the calculated Merkle root
		calculatedMerkleRoot := CalculateMerkleRoot(currentBlock.Transactions)
		if calculatedMerkleRoot != currentBlock.MerkleRoot {
			println("Merkle root mismatch")
			return false
		}

		// Check if the TransactionDateTime of the current block is after the previous block's TransactionDateTime
		if currentBlock.TransactionDateTime.Before(prevBlock.TransactionDateTime) {
			return false
		}
	}
	return true
}

func CalculateHash(block *Block) string {
	data := fmt.Sprintf("%v%d%s%s%s", block.Transactions, block.Nonce, block.PreviousHash, block.MerkleRoot, block.TransactionDateTime)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func CalculateMerkleRoot(transactions []string) string {
	if len(transactions) == 0 {
		return ""
	}
	if len(transactions) == 1 {
		return calculateTransactionHash(transactions[0])
	}

	var merkle []string
	if len(transactions)%2 != 0 {
		transactions = append(transactions, transactions[len(transactions)-1])
	}
	for i := 0; i < len(transactions); i += 2 {
		concatenated := calculateTransactionHash(transactions[i] + transactions[i+1])
		merkle = append(merkle, concatenated)
	}
	return CalculateMerkleRoot(merkle)
}

func calculateTransactionHash(transaction string) string {
	hash := sha256.Sum256([]byte(transaction))
	return fmt.Sprintf("%x", hash)
}

func SetNumberOfTransactionsPerBlock(bc *Blockchain, numTransactions int) {
	bc.NumTransactionsPerBlock = numTransactions
}

func SetBlockHashRangeForBlockCreation(bc *Blockchain, minHash, maxHash string) {
	bc.MinBlockHash = minHash
	bc.MaxBlockHash = maxHash
}

func mineBlock(bc *Blockchain, transactions []string) int {
	for nonce := 0; ; nonce++ {
		data := fmt.Sprintf("%v%d%s%d", transactions, nonce, bc.Blocks[len(bc.Blocks)-1].CurrentHash, time.Now().UnixNano())
		hash := sha256.Sum256([]byte(data))
		hashStr := fmt.Sprintf("%x", hash)
		if hashStr >= bc.MinBlockHash && hashStr <= bc.MaxBlockHash {
			return nonce
		}
	}
}

func main() {
	// Create a new blockchain
	bc := &Blockchain{}

	// Set the number of transactions per block and the block hash range
	SetNumberOfTransactionsPerBlock(bc, 5)
	SetBlockHashRangeForBlockCreation(bc, "0000000000", "00000FFFFF")

	// Create the Genesis block
	genesisTransactions := []string{"Genesis Transaction 1", "Genesis Transaction 2"}
	genesisBlock, _ := NewBlock(bc, genesisTransactions, 0, "", time.Now())
	bc.Blocks = append(bc.Blocks, genesisBlock)

	// Add more transactions and blocks with incremental timestamps
	transactions := []string{"Transaction 1", "Transaction 2", "Transaction 3", "Transaction 4", "Transaction 5"}

	previousBlock := genesisBlock
	for i := 0; i < len(transactions); i += bc.NumTransactionsPerBlock {
		nonce := mineBlock(bc, transactions[i:i+bc.NumTransactionsPerBlock])
		// Use an incremental timestamp for each block
		newBlock, _ := NewBlock(bc, transactions[i:i+bc.NumTransactionsPerBlock], nonce, previousBlock.CurrentHash, previousBlock.TransactionDateTime.Add(time.Second))
		bc.Blocks = append(bc.Blocks, newBlock)
		previousBlock = newBlock
	}

	fmt.Println("--------------------------------------------------------Displaying all blocks in the blockchain--------------------------------------------------------")
	DisplayBlocks(bc)
	fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------")

	// Verify the blockchain
	isValid := VerifyChain(bc)
	if isValid {
		fmt.Println("The blockchain is valid.")
	} else {
		fmt.Println("The blockchain is invalid.")
	}

	// Showcase changing a block's transaction
	targetBlock := bc.Blocks[1]
	fmt.Println("----------------------------------------------------Changing a block's transaction----------------------------------------------------")
	ChangeBlock(targetBlock, "New Transaction", 3)
	fmt.Println("---------------------------------------------------Displaying all blocks in the blockchain after the change---------------------------------------------------")
	DisplayBlocks(bc)
	fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------")

	// Add more blocks with additional transactions
	additionalTransactions := []string{"Transaction 6", "Transaction 7", "Transaction 8", "Transaction 9", "Transaction 10"}
	previousBlock = bc.Blocks[len(bc.Blocks)-1]
	for i := 0; i < len(additionalTransactions); i += bc.NumTransactionsPerBlock {
		nonce := mineBlock(bc, additionalTransactions[i:i+bc.NumTransactionsPerBlock])
		newBlock, _ := NewBlock(bc, additionalTransactions[i:i+bc.NumTransactionsPerBlock], nonce, previousBlock.CurrentHash, previousBlock.TransactionDateTime.Add(time.Second))
		bc.Blocks = append(bc.Blocks, newBlock)
		previousBlock = newBlock
	}
	fmt.Println("---------------------------------------------------Displaying all blocks in the blockchain after adding more blocks---------------------------------------------------")
	DisplayBlocks(bc)
	fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------")

	// Verify the blockchain again
	isValid = VerifyChain(bc)
	if isValid {
		fmt.Println("The blockchain is still valid after changes and additions.")
	} else {
		fmt.Println("The blockchain is invalid after changes and additions.")
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}
