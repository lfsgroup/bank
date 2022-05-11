# Bank

A Go program to look up Australian BSB numbers.  The data is downloaded from [AusPayNet](https://bsb.auspaynet.com.au/), and statically stored in the app.

## Usage

### Go package
You can import `bank` as a Go package into an existing project and use as Go code.

```sh
$ go get -u github.com/timwmillard/bank
```

Then just import ...
```go
import "github.com/timwmillard/bank"
```
and call the `bank.LookupBSB` function ...

```go
branch := bank.LookupBSB("012-023")
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
{"bsb":"012-023","name":"ANZ Wealth Australia Limited","bank":{"code":"ANZ","name":"T\u0026C Town \u0026 Country Bank (a division of Australia \u0026 New Zealand Banking Group Limited)","bsb_numbers":"15"},"address":"347 Kent Street","suburb":"Sydney","state":"NSW","postcode":"2000"}
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

- Institution matching with Branch
  - Once it's matched a branch via `bsb.csv` dataset, it then tried to match the corresponding bank via the `institution.csv` dataset.  It currently matches via the institution code, ie "ANZ" or "CBA", but there are mutiple insitutions with the same name.  Instead we should be matching on the BSB number, ie column 3 in `institution.csv`.

- Loading data will panic
  - The current `loadData` function runs at startup, but will `panic`, if an error occurs.
  - Should handle the error more gracefully.

- Write tests
   - Need to write some tests for the loading and parsing of the data.

## TODO Features

- Write a scrap cli that you can run periodically, to fetch and update to the latest bsb and institution data.

- Better BSB number parser.
  - Current implementation assumes the BSB number matches the value in the data file exactly.
  - Implement a BSB data type, that will accept a number of different formats for the bsb number, ie `123-456` and `123456`.

- Write a cli to look up the bsb numbers.  Usage `$ bsblookup 012-023`.
