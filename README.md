# go-pix

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/fonini/go-pix/pix)
[![Test Status](https://github.com/fonini/go-pix/workflows/tests/badge.svg)](https://github.com/fonini/go-pix/actions?query=workflow%3Atests)
[![codecov](https://codecov.io/gh/fonini/go-pix/branch/main/graph/badge.svg?token=9RNR32U66L&force=true)](https://codecov.io/gh/fonini/go-pix)
[![Go Report Card](https://goreportcard.com/badge/github.com/fonini/go-pix?force=true)](https://goreportcard.com/report/github.com/fonini/go-pix)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

go-pix is a Go library for generating [Pix](https://www.bcb.gov.br/estabilidadefinanceira/pix) transactions using Copy and Paste or QR codes.

## About Pix

![Generated QR code](pix.png?raw=true)

Pix is a system created by the Brazilian Central Bank to allow instant payments. The new payment method allows immediate money transfer, 24 hours a day, 7 days a week, including weekends and holidays.

The address key is a way to identify the user’s account. There are four types of address keys that users can use:

* CPF/CNPJ
* Email address
* Cellphone number
* Random key – a set of random number, letters, and symbols

This key binds the basic information to the user’s complete account information, allowing users to send and receive money using only an address key.

## Usage

```go
import "github.com/fonini/go-pix/pix"
```

### Generating a Copy and Paste code

```go
options := pix.Options{
    Name: "Jonnas Fonini",
    Key: "jonnasfonini@gmail.com",
    City: "Marau",
    Amount: 20.67,
    Description: "Invoice #4",
    TransactionID: "***", // optional
}

copyPaste, err := pix.Pix(options)

if err != nil {
	panic(err)
}

fmt.Println(copyPaste) // will output: "00020126580014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com0210Invoice #4520400005303986540520.675802BR5913Jonnas Fonini6005Marau62410503***50300017BR.GOV.BCB.BRCODE01051.0.06304CF13"
```

### Generating a QR code from a Copy and Paste code

You can use the Copy and Paste code generated above to generate a QR code

```go
options := QRCodeOptions{Size: 256, Content: copyPaste}

qrCode, err := pix.QRCode(options)

if err != nil {
	panic(err)
}
```

The ```qrCode``` is a byte array, containing a graphical representation of the Copy and Paste code in the form of a QR code.

![Generated QR code](qr.png?raw=true)

## Tests

```sh
go test ./pix
```

## License

This open-sourced software is licensed under the [MIT license](https://opensource.org/licenses/MIT).