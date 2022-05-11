package main

import (
	"fmt"
	"os"

	"github.com/timwmillard/bank"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <bsbnumber>\n", os.Args[0])
		os.Exit(1)
	}
	bsb := os.Args[1]
	branch, err := bank.LookupBSB(bsb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lookup BSB %s error: %v\n", bsb, err)
		os.Exit(1)
	}
	fmt.Printf("BSB_NUMBER=%q\n", branch.BSB.String())
	fmt.Printf("BRANK_CODE=%q\n", branch.Bank.Code)
	fmt.Printf("BANK_NAME=%q\n", branch.Bank.Name)
	fmt.Printf("BRANCH_NAME=%q\n", branch.Name)
	fmt.Printf("BRANCH_ADDRESS=%q\n", branch.Address+" "+branch.Suburb+" "+branch.State+" "+branch.Postcode)
}
