package main

import (
	"fmt"

	"github.com/20I-0923/assignment01bca"
)

func main() {
	// Creating a new blockchain
	bc := &assignment01bca.Blockchain{}

	// Adding the genesis block (first block) with an empty previousHash
	block1, err := assignment01bca.NewBlock(bc, "Transaction 1", 123, "")
	if err != nil {
		fmt.Println("Error 1:", err)
		return
	}
	bc.Blocks = append(bc.Blocks, block1)

	// Adding the second block with the CurrentHash of the first block as the previousHash
	block2, err := assignment01bca.NewBlock(bc, "Transaction 2", 456, bc.Blocks[len(bc.Blocks)-1].CurrentHash)
	if err != nil {
		fmt.Println("Error 2:", err)
		return
	}
	bc.Blocks = append(bc.Blocks, block2)

	// Displaying blocks
	assignment01bca.DisplayBlocks(bc)

	// Changing the transaction of the second block
	if len(bc.Blocks) > 1 {
		assignment01bca.ChangeBlock(bc.Blocks[1], "New Transaction 2")
	} else {
		fmt.Println("Blockchain does not have enough blocks to change.")
	}

	assignment01bca.DisplayBlocks(bc)

	// Verify the blockchain
	isValid := assignment01bca.VerifyChain(bc)
	if isValid {
		fmt.Println("Blockchain is valid.")
	} else {
		fmt.Println("Blockchain is not valid.")
	}
}
