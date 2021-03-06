# Bank

A Go program to look up Australian BSB numbers.  The data is downloaded from [AusPayNet](https://bsb.auspaynet.com.au/), and statically stored in the app.

## Usage

### CLI
Install the `bsblookup` cli, to use it as a command.
```sh
$ go install github.com/lfsgroup/bank/cmd/bsblookup@latest
```

Then run from the command line prompt:
```sh
$ bsblookup 012-023
BSB_NUMBER="012-023"
BANK_CODE="ANZ"
BANK_NAME="Australia & New Zealand Banking Group Limited"
BRANCH_NAME="ANZ Wealth Australia Limited"
BRANCH_ADDRESS="347 Kent Street Sydney NSW 2000"
BRANCH_PAYMENTS_FLAGS="PEH"
```

### Go package
You can import `bank` as a Go package, and use in an existing Go project.

```sh
$ go get -u github.com/lfsgroup/bank
```

Then just import ...
```go
import "github.com/lfsgroup/bank"
```
and call the `bank.LookupBSB` function ...

```go
branch, err := bank.LookupBSB("012-023")
if err != nil {
   log.Fatalln("Lookup BSB error:", err)
}
fmt.Println("Bank name for 012-023 is", branch.Bank.Name)
```

### HTTP Service

`bankd` is a standalone webserver that can be used to request a bank's BSB number over HTTP.

Build the app ...
```sh
$ go build ./cmd/bankd
```

Run the app ...
```sh
$ PORT=4000 ./bankd
2022/05/11 12:40:50 Server starting on port 4000
```

```sh
$ curl http://localhost:4000/bsb/012-023
{"bsb":"012-023","name":"ANZ Wealth Australia Limited","bank":{"code":"ANZ","name":"Australia \u0026 New Zealand Banking Group Limited","bsb_numbers":"1"},"bank_code":"ANZ","address":"347 Kent Street","suburb":"Sydney","state":"NSW","postcode":"2000","payments_flags":"PEH"}
```

### AWS Lambda

You can also deploy the `bank` as a standalone AWS Lambda service.

Build the lambda
```sh
$ make build
```
This will build a file call `build/bsblookup.zip`.  You can either upload this file directly in the AWS Lambda console or run the following command using aws cli:
```sh
$ aws lambda update-function-code --function-name bsblookup --zip-file fileb://$PWD/build/bsblookup.zip
```
This assumed you created a Lambda function called `bsblookup`.

## Bugs / Issues

- Write tests
   - Need to write some tests for the loading and parsing of the data.
