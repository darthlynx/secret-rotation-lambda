package generator

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"

	"github.com/darthlynx/secret-rotation-lambda/internal/config"
	"github.com/darthlynx/secret-rotation-lambda/internal/models"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	specialChars   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	ambiguousChars = "0Ol1"
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

func (g *SecretGenerator) Generate(opts models.GeneratorOptions) (string, error) {
	if err := config.ValidateGeneratorOptions(opts); err != nil {
		return "", err
	}

	charset := buildCharset(opts)
	if len(charset) == 0 {
		return "", errors.New("no character types selected")
	}

	secret := make([]byte, opts.Length)
	for i := range secret {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		secret[i] = charset[idx.Int64()]
	}
	return string(secret), nil
}

func buildCharset(opts models.GeneratorOptions) string {
	var charset strings.Builder

	if opts.IncludeLowercase {
		charset.WriteString(lowercaseChars)
	}
	if opts.IncludeUppercase {
		charset.WriteString(uppercaseChars)
	}
	if opts.IncludeDigits {
		charset.WriteString(digitChars)
	}
	if opts.IncludeSpecialChars {
		charset.WriteString(specialChars)
	}

	result := charset.String()

	if opts.ExcludeAmbiguous {
		result = removeAmbiguousChars(result)
	}

	return result
}

func removeAmbiguousChars(charset string) string {
	var result strings.Builder
	for _, ch := range charset {
		if !strings.ContainsRune(ambiguousChars, ch) {
			result.WriteRune(ch)
		}
	}
	return result.String()
}
