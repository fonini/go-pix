package pix

import (
	"bytes"
	"errors"
	"image"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// test that input matches the value we want. If not, report an error on t.
func testValue(t *testing.T, input Options, want string) {
	v, err := Pix(input)

	if err != nil {
		t.Errorf("Pix(%v) returned an error: %v", input, err)
	}

	if diff := cmp.Diff(want, v); diff != "" {
		t.Errorf("Pix(%v) mismatch:\n%s", input, diff)
	}
}

func testError(t *testing.T, input Options, want error) {
	v, err := Pix(input)

	if err == nil {
		t.Errorf("Expected error for input %v but received %v", input, v)
	}

	if err.Error() != want.Error() {
		t.Errorf("Unexpected error for input %v: %v (expected %v)", input, err, want)
	}
}

func TestValues_FullOptions(t *testing.T) {
	tests := []struct {
		input Options
		want  string
	}{
		{Options{
			Name:        "Jonnas Fonini",
			Key:         "jonnasfonini@gmail.com",
			City:        "Marau",
			Amount:      20.67,
			Description: "Invoice #4",
		}, "00020126580014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com0210Invoice #4520400005303986540520.675802BR5913Jonnas Fonini6005Marau62410503***50300017BR.GOV.BCB.BRCODE01051.0.06304CF13"},
		{Options{
			Name:        "Jonnas Fonini",
			Key:         "+5554000000000",
			City:        "Porto Alegre",
			Amount:      5.50,
			Description: "Lunch money",
		}, "00020126510014BR.GOV.BCB.PIX0114+55540000000000211Lunch money52040000530398654045.505802BR5913Jonnas Fonini6012Porto Alegre62410503***50300017BR.GOV.BCB.BRCODE01051.0.063044766"},
	}

	for _, tt := range tests {
		testValue(t, tt.input, tt.want)
	}
}

func TestValues_TransactionID(t *testing.T) {
	tests := []struct {
		input Options
		want  string
	}{
		{Options{
			Name:          "Jonnas Fonini",
			Key:           "jonnasfonini@gmail.com",
			City:          "Marau",
			Amount:        20.67,
			Description:   "Invoice #4",
			TransactionID: "99999",
		}, "00020126580014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com0210Invoice #4520400005303986540520.675802BR5913Jonnas Fonini6005Marau624305059999950300017BR.GOV.BCB.BRCODE01051.0.0630401F3"},
	}

	for _, tt := range tests {
		testValue(t, tt.input, tt.want)
	}
}

func TestValues_WithoutAmount(t *testing.T) {
	tests := []struct {
		input Options
		want  string
	}{
		{Options{
			Name: "Jonnas Fonini",
			Key:  "jonnasfonini@gmail.com",
			City: "Marau",
		}, "00020126480014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com020052040000530398654040.005802BR5913Jonnas Fonini6005Marau62410503***50300017BR.GOV.BCB.BRCODE01051.0.06304CC71"},
	}

	for _, tt := range tests {
		testValue(t, tt.input, tt.want)
	}
}

func TestValues_Errors(t *testing.T) {
	tests := []struct {
		input Options
		want  error
	}{
		{Options{}, errors.New("key must not be empty")},
		{Options{
			Key:  "jonnasfonini@gmail.com",
			Name: "Receiver long name to cause error",
			City: "Marau",
		}, errors.New("name must be at least 25 characters long")},
		{Options{
			Key:  "jonnasfonini@gmail.com",
			Name: "Jonnas",
			City: "Receiver city long name",
		}, errors.New("city must be at least 15 characters long")},
		{Options{
			Name: "Jonnas",
			Key:  "jonnasfonini@gmail.com",
		}, errors.New("city must not be empty")},
		{Options{
			City: "Marau",
			Key:  "jonnasfonini@gmail.com",
		}, errors.New("name must not be empty")},
	}

	for _, tt := range tests {
		testError(t, tt.input, tt.want)
	}
}

// Generate a QR Code from a Pix Copy and Paste string and decode the result
func TestQrCodeContent(t *testing.T) {
	str := "00020126580014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com0210Invoice #4520400005303986540520.675802BR5913Jonnas Fonini6005Marau62410503***50300017BR.GOV.BCB.BRCODE01051.0.06304CF13"

	options := QRCodeOptions{Content: str}

	qr, err := QRCode(options)

	if err != nil {
		t.Errorf("QRCode(%v) returned an error: %v", options, err)
	}

	img, _, err := image.Decode(bytes.NewReader(qr))

	if err != nil {
		t.Errorf("image.Decode returned an error: %v", err)
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)

	if err != nil {
		t.Errorf("gozxing.NewBinaryBitmapFromImage returned an error: %v", err)
	}

	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)

	if err != nil {
		t.Errorf("qrReader.Decode returned an error: %v", err)
	}

	if diff := cmp.Diff(result.String(), str); diff != "" {
		t.Errorf("QR Code content is unexpected:\n%s", diff)
	}
}
