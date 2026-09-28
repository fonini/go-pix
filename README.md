# go-pix

[![GoDoc](https://pkg.go.dev/badge/github.com/fonini/go-pix)](https://pkg.go.dev/github.com/fonini/go-pix/pix)
[![Test Status](https://github.com/fonini/go-pix/workflows/tests/badge.svg)](https://github.com/fonini/go-pix/actions?query=workflow%3Atests)
[![golangci-lint](https://github.com/fonini/go-pix/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/fonini/go-pix/actions/workflows/golangci-lint.yml)
[![codecov](https://codecov.io/gh/fonini/go-pix/branch/main/graph/badge.svg?token=9RNR32U66L&force=true)](https://codecov.io/gh/fonini/go-pix)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

go-pix is a Go library for generating and reading [Pix](https://www.bcb.gov.br/estabilidadefinanceira/pix) Copy-and-Paste payloads and PNG QR codes.

## About Pix

![Generated QR code](pix.png?raw=true)

Pix is a system created by the Brazilian Central Bank to allow instant payments. The new payment method allows immediate money transfer, 24 hours a day, 7 days a week, including weekends and holidays.

The address key is a way to identify the user’s account. There are four types of address keys that users can use:

* CPF/CNPJ
* Email address
* Cellphone number
* Random key – a set of random number, letters, and symbols

This key binds the basic information to the user’s complete account information, allowing users to send and receive money using only an address key.

## Installation

```sh
go get github.com/fonini/go-pix@v1.1.2
```

## Quick start

```go
package main

import (
	"fmt"
	"log"

	"github.com/fonini/go-pix/pix"
)

func main() {
	payload, err := pix.Pix(pix.Options{
		Key:           "merchant@example.com",
		Name:          "Example Store",
		City:          "Sao Paulo",
		Amount:        20.67,
		Description:   "Order #123",
		TransactionID: "ORDER123",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(payload)
}
```

## Options

| Field | Required | Rules |
| --- | :---: | --- |
| `Key` | Yes | Pix key: CPF/CNPJ, email, phone number, or random key. |
| `Name` | Yes | Receiver name; at most 25 characters. |
| `City` | Yes | Receiver city; at most 15 characters. |
| `Amount` | No | Amount in BRL. `0` means an open-value Pix. |
| `Description` | No | Description shown with the payment. |
| `TransactionID` | No | At most 25 alphanumeric characters. Leave it empty to use the default identifier. |

## Generating a Copy-and-Paste payload

```go
options := pix.Options{
    Name: "Jonnas Fonini",
    Key: "jonnasfonini@gmail.com",
    City: "Marau",
    Amount: 20.67, // optional
    Description: "Invoice #4", // optional
}

copyPaste, err := pix.Pix(options)

if err != nil {
    log.Fatal(err)
}

fmt.Println(copyPaste) // will output: "00020126580014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com0210Invoice #4520400005303986540520.675802BR5913Jonnas Fonini6005Marau62410503***50300017BR.GOV.BCB.BRCODE01051.0.06304CF13"

optionsFromCode, err := pix.ReadPix(copyPaste)
if err != nil {
    log.Fatal(err)
}

fmt.Println(optionsFromCode)
```

`ReadPix` returns an error for malformed payloads.

## Generating an open-value Pix

Omit `Amount`, or set it to `0`, to generate a payload where the payer chooses the amount.

```go
payload, err := pix.Pix(pix.Options{
    Key:  "merchant@example.com",
    Name: "Example Store",
    City: "Sao Paulo",
})
if err != nil {
    log.Fatal(err)
}
```

## Generating a QR code from a Copy-and-Paste payload

You can use the Copy and Paste code generated above to generate a QR code

```go
func saveQRCode(copyPaste string) error {
    options := pix.QRCodeOptions{Size: 256, Content: copyPaste}

    qrCode, err := pix.QRCode(options)
    if err != nil {
        return err
    }

    return os.WriteFile("pix.png", qrCode, 0o644)
}
```

The example above requires `import "os"`.

`QRCode` returns PNG bytes containing a graphical representation of the Copy-and-Paste payload.

![Generated QR code](qr.png?raw=true)

## Banks tested
* Caixa Econômica Federal
* Nubank
* PicPay
* PagSeguro
* Itaú
* Mercado Pago

## Tests

```sh
go test ./...
```

## License

This open-sourced software is licensed under the [MIT license](https://opensource.org/licenses/MIT).
