package generator

import (
	"fmt"

	"github.com/darthlynx/secret-rotation-lambda/internal/models"
	"github.com/sethvargo/go-password/password"
)

const (
	MinSecretLength  = 8
	MinNumberDigits  = 1
	MinNumberSpecial = 1
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
func (g *SecretGenerator) Generate(opts models.GeneratorOptions) (string, error) {
	if opts.Length < MinSecretLength {
		return "", fmt.Errorf("length must be at least %d", MinSecretLength)
	}
	numDigits := 0
	if opts.IncludeDigits  {
		if opts.MinNumberDigits > 0 {
			numDigits = opts.MinNumberDigits
		} else {
			numDigits = MinNumberDigits
		}
	}
	numSymbols := 0
	if opts.IncludeSpecialChars {
		if opts.MinNumberSpecial > 0 {
			numSymbols = opts.MinNumberSpecial
		} else {
			numSymbols = MinNumberSpecial
		}
	}
	secret, err := password.Generate(opts.Length, numDigits, numSymbols, opts.IncludeUppercase, true) // allow characters repeat
	if err != nil {
		return "", err
	}
	return string(secret), nil
}
