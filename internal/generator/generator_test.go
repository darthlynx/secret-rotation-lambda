package generator

import (
	"testing"

	"unicode"

	"github.com/darthlynx/secret-rotation-lambda/internal/models"
)

func TestGenerate(t *testing.T) {
	gen := New()

	tests := []struct {
		name    string
		opts    models.GeneratorOptions
		wantErr bool
	}{
		{
			name: "valid with all character types",
			opts: models.GeneratorOptions{
				Length:              16,
				IncludeDigits:       true,
				IncludeUppercase:    true,
				IncludeSpecialChars: true,
			},
			wantErr: false,
		},
		{
			name: "lowercase only",
			opts: models.GeneratorOptions{
				Length: 12,
			},
			wantErr: false,
		},
		{
			name: "invalid - zero length",
			opts: models.GeneratorOptions{
				Length: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret, err := gen.Generate(tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Generate() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Generate() unexpected error: %v", err)
				return
			}

			if len(secret) != tt.opts.Length {
				t.Errorf("Generate() length = %d, want %d", len(secret), tt.opts.Length)
			}
		})
	}
}

func TestGenerateDigitsLength(t *testing.T) {
	gen := New()

	tests := []struct {
		name    string
		opts    models.GeneratorOptions
		wantLen int
	}{
		{
			name: "check digit length is 0 when not included",
			opts: models.GeneratorOptions{
				Length:        16,
				IncludeDigits: false,
			},
			wantLen: 0,
		},
		{
			name: "check digits length is 1 of total length if MinNumberDigits not set",
			opts: models.GeneratorOptions{
				Length:        16,
				IncludeDigits: true,
			},
			wantLen: 1,
		},
		{
			name: "check digits lengths is set by MinNumberDigits",
			opts: models.GeneratorOptions{
				Length:          24,
				IncludeDigits:   true,
				MinNumberDigits: 20,
			},
			wantLen: 20,
		},
		{
			name: "check special chars is set to 0 when not included",
			opts: models.GeneratorOptions{
				Length:              10,
				IncludeSpecialChars: false,
			},
			wantLen: 0,
		},
		{
			name: "check special chars is set to 1 when included and MinNumberSpecial not set",
			opts: models.GeneratorOptions{
				Length:              10,
				IncludeSpecialChars: true,
			},
			wantLen: 1,
		},
		{
			name: "check special chars is set by MinNumberSpecial",
			opts: models.GeneratorOptions{
				Length:              10,
				IncludeSpecialChars: true,
				MinNumberSpecial:    7,
			},
			wantLen: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret, err := gen.Generate(tt.opts)

			if err != nil {
				t.Errorf("Generate() unexpected error: %v", err)
				return
			}

			charCount := 0
			specialCharCount := 0
			for _, ch := range secret {
				if unicode.IsDigit(ch) {
					charCount++
				} else if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
					specialCharCount++
				}
			}

			if tt.opts.IncludeDigits && !tt.opts.IncludeSpecialChars {
				if charCount != tt.wantLen {
					t.Errorf("Generate() digits count = %d, want %d", charCount, tt.wantLen)
				}
			}
			if tt.opts.IncludeSpecialChars && !tt.opts.IncludeDigits {
				if specialCharCount != tt.wantLen {
					t.Errorf("Generate() special chars count = %d, want %d", specialCharCount, tt.wantLen)
				}
			}
			if tt.opts.IncludeSpecialChars && tt.opts.IncludeDigits {
				if (charCount + specialCharCount) != tt.wantLen {
					t.Errorf("Generate() digits + special chars count = %d, want %d", charCount+specialCharCount, tt.wantLen)
				}
			}
		})
	}
}

func TestGenerateUniqueness(t *testing.T) {
	gen := New()
	opts := models.GeneratorOptions{
		Length:              32,
		IncludeDigits:       true,
		IncludeUppercase:    true,
		IncludeSpecialChars: true,
	}

	secrets := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		secret, err := gen.Generate(opts)
		if err != nil {
			t.Fatalf("Generate() error: %v", err)
		}
		if secrets[secret] {
			t.Errorf("Generate() produced duplicate secret")
		}
		secrets[secret] = true
	}
}
