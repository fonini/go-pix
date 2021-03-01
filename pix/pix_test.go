package pix

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// test that input matches the value we want. If not, report an error on t.
func testValue(t *testing.T, input Options, want string) {
	v, err := Pix(input)

	if err != nil {
		t.Errorf("Pix(%v) returned error: %v", input, err)
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
		}, "00020126580014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com0210Invoice #4520400005303986540520.675802BR5913Jonnas Fonini6005Marau624305059999950300017BR.GOV.BCB.BRCODE01051.0.063041F3"},
	}

	for _, tt := range tests {
		testValue(t, tt.input, tt.want)
	}
}

func TestValues_WithoutName(t *testing.T) {
	tests := []struct {
		input Options
		want  string
	}{
		{Options{
			Key:    "jonnasfonini@gmail.com",
			City:   "Marau",
			Amount: 20.67,
		}, "00020126480014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com0200520400005303986540520.675802BR6005Marau62410503***50300017BR.GOV.BCB.BRCODE01051.0.06304948B"},
	}

	for _, tt := range tests {
		testValue(t, tt.input, tt.want)
	}
}

func TestValues_WithoutCity(t *testing.T) {
	tests := []struct {
		input Options
		want  string
	}{
		{Options{
			Name:   "Jonnas Fonini",
			Key:    "jonnasfonini@gmail.com",
			Amount: 20.67,
		}, "00020126480014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com0200520400005303986540520.675802BR5913Jonnas Fonini62410503***50300017BR.GOV.BCB.BRCODE01051.0.06304BF73"},
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

func TestValues_OnlyAmountAndAccount(t *testing.T) {
	tests := []struct {
		input Options
		want  string
	}{
		{Options{
			Key:    "jonnasfonini@gmail.com",
			Amount: 5.50,
		}, "00020126480014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com020052040000530398654045.505802BR62410503***50300017BR.GOV.BCB.BRCODE01051.0.06304E0EE"},
	}

	for _, tt := range tests {
		testValue(t, tt.input, tt.want)
	}
}

func TestValues_OnlyAccount(t *testing.T) {
	tests := []struct {
		input Options
		want  string
	}{
		{Options{
			Key: "jonnasfonini@gmail.com",
		}, "00020126480014BR.GOV.BCB.PIX0122jonnasfonini@gmail.com020052040000530398654040.005802BR62410503***50300017BR.GOV.BCB.BRCODE01051.0.0630430B9"},
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
		}, errors.New("name must be at least 25 characters long")},
		{Options{
			Key:  "jonnasfonini@gmail.com",
			City: "Receiver city long name",
		}, errors.New("city must be at least 15 characters long")},
	}

	for _, tt := range tests {
		testError(t, tt.input, tt.want)
	}
}
