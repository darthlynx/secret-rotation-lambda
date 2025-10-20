package generator

import (
	"strings"
	"testing"

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
				IncludeLowercase:    true,
				IncludeSpecialChars: true,
			},
			wantErr: false,
		},
		{
			name: "lowercase only",
			opts: models.GeneratorOptions{
				Length:           12,
				IncludeLowercase: true,
			},
			wantErr: false,
		},
		{
			name: "exclude ambiguous characters",
			opts: models.GeneratorOptions{
				Length:           20,
				IncludeDigits:    true,
				IncludeUppercase: true,
				ExcludeAmbiguous: true,
			},
			wantErr: false,
		},
		{
			name: "invalid - zero length",
			opts: models.GeneratorOptions{
				Length:           0,
				IncludeLowercase: true,
			},
			wantErr: true,
		},
		{
			name: "invalid - no character types",
			opts: models.GeneratorOptions{
				Length: 10,
			},
			wantErr: true,
		},
		{
			name: "invalid - length too large",
			opts: models.GeneratorOptions{
				Length:           5000,
				IncludeLowercase: true,
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

			if tt.opts.ExcludeAmbiguous {
				for _, ch := range secret {
					if strings.ContainsRune(ambiguousChars, ch) {
						t.Errorf("Generate() contains ambiguous character: %c", ch)
					}
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
		IncludeLowercase:    true,
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
