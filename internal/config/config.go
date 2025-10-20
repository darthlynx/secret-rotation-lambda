package config

import (
	"errors"

	"github.com/darthlynx/secret-rotation-lambda/internal/models"
)

func ValidateRotationRequest(req models.RotationRequest) error {
	if err := validateSecretARN(req.SecretARN); err != nil {
		return err
	}

	if err := validateSecretType(req.SecretType); err != nil {
		return err
	}

	if err := ValidateGeneratorOptions(req.GeneratorOpts); err != nil {
		return err
	}

	if req.SecretType == models.SecretTypeKeyValue || req.SecretType == models.SecretTypeJSON {
		if err := validateKeyValueConfig(req.KeyValueConfig); err != nil {
			return err
		}
	}

	return nil
}

func validateSecretARN(arn string) error {
	if arn == "" {
		return errors.New("secret_arn cannot be empty")
	}

	// Basic length check for ARN format
	if len(arn) < 20 {
		return errors.New("secret_arn is too short to be valid")
	}

	return nil
}

func validateSecretType(secretType models.SecretType) error {
	switch secretType {
	case models.SecretTypePlaintext, models.SecretTypeKeyValue, models.SecretTypeJSON:
		return nil
	default:
		return errors.New("invalid secret_type")
	}
}

func ValidateGeneratorOptions(opts models.GeneratorOptions) error {
	if opts.Length < 8 || opts.Length > 2048 {
		return errors.New("length must be between 8 and 2048")
	}

	hasCharType := opts.IncludeDigits || opts.IncludeUppercase || opts.IncludeLowercase || opts.IncludeSpecialChars
	if !hasCharType {
		return errors.New("at least one character type must be included")
	}
	return nil
}

func validateKeyValueConfig(cfg *models.KeyValueConfig) error {
	if cfg == nil {
		return nil
	}
	// Empty KeysToRotate means rotate all keys, which is valid
	return nil
}
