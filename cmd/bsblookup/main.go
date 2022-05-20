package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/timwmillard/bank"
)

func main() {
	j := flag.Bool("j", false, "Output JSON format")
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	bsb := flag.Args()[0]
	branch, err := bank.LookupBSB(bsb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lookup BSB %s error: %v\n", bsb, err)
		os.Exit(1)
	}

	if *j {
		printJSON(branch)
	} else {
		printStandard(branch)
	}
}

func printStandard(branch bank.Branch) {
	fmt.Printf("BSB_NUMBER=%q\n", branch.BSB.String())
	fmt.Printf("BANK_CODE=%q\n", branch.Bank.Code)
	fmt.Printf("BANK_NAME=%q\n", branch.Bank.Name)
	fmt.Printf("BRANCH_NAME=%q\n", branch.Name)
	fmt.Printf("BRANCH_ADDRESS=%q\n", branch.Address+" "+branch.Suburb+" "+branch.State+" "+branch.Postcode)
}

func printJSON(branch bank.Branch) {
	err := json.NewEncoder(os.Stdout).Encode(branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encoding JSON error: %v\n", err)
	}
}
