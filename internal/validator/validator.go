package validator

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

func validateKeyValueConfig(cfg *models.KeyValueConfig) error {
	if cfg == nil {
		return nil
	}
	// Empty KeysToRotate means rotate all keys, which is valid
	return nil
}
