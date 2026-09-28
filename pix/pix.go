// Package pix generates Brazilian Pix Copy and Paste or QR Codes
//
// As a simple example:
//
//	options := pix.Options{
//		Name: "Jonnas Fonini",
//		Key: "jonnasfonini@gmail.com",
//		City: "Marau",
//		Amount: 20.67, // optional
//		Description: "Invoice #4", // optional
//		TransactionID: "***", // optional
//	}
//
//	copyPaste, err := pix.Pix(options)
//
//	if err != nil {
//		fmt.Println("could not generate Pix:", err)
//		return
//	}
//
//	fmt.Println(copyPaste) // will output: "00020126580014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com0210Invoice #4520400005303986540520.675802BR5913Jonnas Fonini6005Marau62410503***50300017BR.GOV.BCB.BRCODE01051.0.06304CF13"
package pix

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"unicode/utf8"

	"github.com/r10r/crc16"
	"github.com/skip2/go-qrcode"
)

// Options is a configuration struct.
type Options struct {
	// Pix Key (CPF/CNPJ, Email, Cellphone or Random Key)
	Key string
	// Receiver name
	Name string
	// Receiver city
	City string
	// Transaction amount
	Amount float64
	// Transaction description
	Description string
	// Transaction ID
	TransactionID string
}

// QRCodeOptions is a configuration struct.
type QRCodeOptions struct {
	// QR Code content
	Content string
	// Default: 256
	Size int
}

type intMap map[int]interface{}

// Pix generates a Copy and Paste Pix code
func Pix(options Options) (string, error) {
	if err := validateData(options); err != nil {
		return "", err
	}

	data := buildDataMap(options)
	str := parseData(data)

	// Add the CRC at the end
	str += "6304"

	crc, err := calculateCRC16(str)

	if err != nil {
		return "", err
	}

	str += crc

	return str, nil
}

// QRCode returns a graphical representation of the Copy and Paste code in a QR Code form.
func QRCode(options QRCodeOptions) ([]byte, error) {
	if options.Size == 0 {
		options.Size = 256
	}

	bytes, err := qrcode.Encode(options.Content, qrcode.Medium, options.Size)

	return bytes, err
}

func validateData(options Options) error {
	if options.Key == "" {
		return errors.New("key must not be empty")
	}

	if options.Name == "" {
		return errors.New("name must not be empty")
	}

	if options.City == "" {
		return errors.New("city must not be empty")
	}

	if utf8.RuneCountInString(options.Name) > 25 {
		return errors.New("name must be at most 25 characters long")
	}

	if utf8.RuneCountInString(options.City) > 15 {
		return errors.New("city must be at most 15 characters long")
	}

	// Validate Transaction ID when provided
	if options.TransactionID != "" {
		// Max length 25 characters
		if utf8.RuneCountInString(options.TransactionID) > 25 {
			return errors.New("transaction id must be at most 25 characters long")
		}
		// Only alphanumeric characters
		for _, r := range options.TransactionID {
			if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
				return errors.New("transaction id must be alphanumeric (letters and numbers only)")
			}
		}
	}

	return nil
}

func buildDataMap(options Options) intMap {
	data := make(intMap)

	// Payload Format Indicator
	data[0] = "01"

	// Merchant Account Information
	data[26] = intMap{0: "BR.GOV.BCB.PIX", 1: options.Key, 2: options.Description}

	// Merchant Category Code
	data[52] = "0000"

	// Transaction Currency - Brazilian Real - ISO4217
	data[53] = "986"

	// Transaction Amount. Omit it when no amount is defined, so the payer can enter it.
	if options.Amount != 0 {
		data[54] = options.Amount
	}

	// Country Code - ISO3166-1 alpha 2
	data[58] = "BR"

	// Merchant Name. 25 characters maximum
	data[59] = options.Name

	// Merchant City. 15 characters maximum
	data[60] = options.City

	// Transaction ID
	data[62] = intMap{5: "***", 50: intMap{0: "BR.GOV.BCB.BRCODE", 1: "1.0.0"}}

	if options.TransactionID != "" {
		data[62].(intMap)[5] = options.TransactionID
	}

	return data
}

func parseData(data intMap) string {
	var str string

	keys := sortKeys(data)

	for _, k := range keys {
		v := reflect.ValueOf(data[k])

		switch v.Kind() {
		case reflect.String:
			value := data[k].(string)
			str += fmt.Sprintf("%02d%02d%s", k, len(value), value)
		case reflect.Float64:
			value := strconv.FormatFloat(v.Float(), 'f', 2, 64)

			str += fmt.Sprintf("%02d%02d%s", k, len(value), value)
		case reflect.Map:
			// If the element is another map, do a recursive call
			content := parseData(data[k].(intMap))

			str += fmt.Sprintf("%02d%02d%s", k, len(content), content)
		}
	}

	return str
}

// ReadPix generates an Options struct using a copyPaste PIX code
func ReadPix(copyPaste string) (Options, error) {
	data, err := buildUsingGuideMap(copyPaste, buildDataMap(Options{}))
	if err != nil {
		return Options{}, err
	}

	return readDataMap(data)
}

func readDataMap(data intMap) (op Options, err error) {
	keyMap, ok := data[26].(intMap)
	if !ok {
		return op, errors.New("invalid Pix code: missing merchant account information")
	}
	key, ok := keyMap[1].(string)
	if !ok {
		return op, errors.New("invalid Pix code: missing Pix key")
	}
	description, ok := keyMap[2].(string)
	if !ok {
		return op, errors.New("invalid Pix code: missing transaction description")
	}
	name, ok := data[59].(string)
	if !ok {
		return op, errors.New("invalid Pix code: missing merchant name")
	}
	city, ok := data[60].(string)
	if !ok {
		return op, errors.New("invalid Pix code: missing merchant city")
	}
	txMap, ok := data[62].(intMap)
	if !ok {
		return op, errors.New("invalid Pix code: missing additional data")
	}
	transactionID, ok := txMap[5].(string)
	if !ok {
		return op, errors.New("invalid Pix code: missing transaction ID")
	}
	if transactionID == "***" {
		transactionID = ""
	}

	amount, _ := data[54].(float64)

	op = Options{
		Key:           key,
		Description:   description,
		Amount:        amount,
		Name:          name,
		City:          city,
		TransactionID: transactionID,
	}

	return op, err
}

func buildUsingGuideMap(copyPaste string, guide intMap) (intMap, error) {
	data := make(intMap)

	k := 0
	for k < len(copyPaste) {
		if len(copyPaste)-k < 4 {
			return nil, errors.New("invalid Pix code: truncated TLV header")
		}

		index, err := strconv.Atoi(copyPaste[k : k+2])
		if err != nil {
			return nil, errors.New("invalid Pix code: invalid TLV tag")
		}
		k += 2

		length, err := strconv.Atoi(copyPaste[k : k+2])
		if err != nil {
			return nil, errors.New("invalid Pix code: invalid TLV length")
		}
		k += 2
		if length > len(copyPaste)-k {
			return nil, errors.New("invalid Pix code: truncated TLV value")
		}

		value := copyPaste[k : k+length]
		k += length

		if index == 54 {
			data[index], _ = strconv.ParseFloat(value, 64)
			continue
		}

		v := reflect.ValueOf(guide[index])
		switch v.Kind() {
		case reflect.Map:
			m := guide[index].(intMap)
			decoded, err := buildUsingGuideMap(value, m)
			if err != nil {
				return nil, err
			}
			data[index] = decoded
		case reflect.String:
			data[index] = value
		case reflect.Float64:
			data[index], _ = strconv.ParseFloat(value, 64)
		}
	}

	return data, nil
}

func sortKeys(data intMap) []int {
	keys := make([]int, len(data))
	i := 0

	for k := range data {
		keys[i] = k
		i++
	}

	sort.Ints(keys)

	return keys
}

func calculateCRC16(str string) (string, error) {
	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)

	h := crc16.New(table)
	_, err := h.Write([]byte(str))

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%04X", h.Sum16()), nil
}
