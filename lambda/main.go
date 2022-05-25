package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lfsgroup/bank"
)

func LookupBSB(ctx context.Context, bsb string) (bank.Branch, error) {
	log.Println("Lookup BSB: ", bsb)
	return bank.LookupBSB(bsb)
}

func main() {
	log.Println("Bank BSB Lookup API: Starting")

	lambda.Start(LookupBSB)
}
