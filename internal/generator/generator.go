package generator

import (
	"errors"

	"github.com/darthlynx/secret-rotation-lambda/internal/models"
	"github.com/sethvargo/go-password/password"
)

// Generator defines the interface for secret generation.
type Generator interface {
	Generate(opts models.GeneratorOptions) (string, error)
}

// SecretGenerator implements the Generator interface.
type SecretGenerator struct{}

// New creates a new instance of SecretGenerator.
func New() *SecretGenerator {
	return &SecretGenerator{}
}

// Generate creates a new secret based on the provided options
//
// Number of digits and symbols are calculated as 1/4 of the total length if included.
func (g *SecretGenerator) Generate(opts models.GeneratorOptions) (string, error) {
	if opts.Length < 8 {
		return "", errors.New("length must be at least 8")
	}
	numDigits := 0
	if opts.IncludeDigits {
		numDigits = opts.Length / 4
	}
	numSymbols := 0
	if opts.IncludeSpecialChars {
		numSymbols = opts.Length / 4
	}
	secret, err := password.Generate(opts.Length, numDigits, numSymbols, opts.IncludeUppercase, true) // allow characters repeat
	if err != nil {
		return "", err
	}
	return string(secret), nil
}
