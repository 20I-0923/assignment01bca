// main.go

package main

import (
	"fmt"

	"github.com/20I-0923/assignment01bca/assignment01bca"
)

func main() {
	// Create a new blockchain
	bc := &assignment01bca.Blockchain{}

	// Add blocks to the blockchain
	bc.Blocks = append(bc.Blocks, assignment01bca.NewBlock("Transaction 1", 123, ""))
	bc.Blocks = append(bc.Blocks, assignment01bca.NewBlock("Transaction 2", 456, bc.Blocks[len(bc.Blocks)-1].CurrentHash))

	// Display the blocks
	assignment01bca.DisplayBlocks(bc)

	// Change the transaction of the second block
	assignment01bca.ChangeBlock(bc.Blocks[1], "New Transaction 2")

	// Verify the blockchain
	isValid := assignment01bca.VerifyChain(bc)
	if isValid {
		fmt.Println("Blockchain is valid.")
	} else {
		fmt.Println("Blockchain is not valid.")
	}
}
