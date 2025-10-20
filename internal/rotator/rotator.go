package rotator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/darthlynx/secret-rotation-lambda/internal/validator"
	"github.com/darthlynx/secret-rotation-lambda/internal/generator"
	"github.com/darthlynx/secret-rotation-lambda/internal/models"
	"github.com/darthlynx/secret-rotation-lambda/internal/secretsmanager"
)

// Rotator handles secret rotation logic.
type Rotator struct {
	smClient secretsmanager.Client
	gen      generator.Generator
}

func New(smClient secretsmanager.Client, gen generator.Generator) *Rotator {
	return &Rotator{
		smClient: smClient,
		gen:      gen,
	}
}

// RotateSecret performs the secret rotation based on the request.
func (r *Rotator) RotateSecret(ctx context.Context, req models.RotationRequest) (*models.RotationResponse, error) {
	if err := validator.ValidateRotationRequest(req); err != nil {
		return &models.RotationResponse{
			Success:   false,
			SecretARN: req.SecretARN,
			ErrorMsg:  err.Error(),
		}, err
	}

	var newSecretValue string
	var err error

	switch req.SecretType {
	case models.SecretTypePlaintext:
		newSecretValue, err = r.rotatePlaintext(ctx, req)
	case models.SecretTypeKeyValue, models.SecretTypeJSON:
		newSecretValue, err = r.rotateKeyValue(ctx, req)
	default:
		err = fmt.Errorf("unsupported secret type: %s", req.SecretType)
	}

	if err != nil {
		return &models.RotationResponse{
			Success:   false,
			SecretARN: req.SecretARN,
			ErrorMsg:  err.Error(),
		}, err
	}

	versionID, err := r.smClient.PutSecretValue(ctx, req.SecretARN, newSecretValue)
	if err != nil {
		return &models.RotationResponse{
			Success:   false,
			SecretARN: req.SecretARN,
			ErrorMsg:  fmt.Sprintf("failed to update secret: %v", err),
		}, err
	}

	return &models.RotationResponse{
		Success:   true,
		SecretARN: req.SecretARN,
		VersionID: versionID,
	}, nil
}

func (r *Rotator) rotatePlaintext(ctx context.Context, req models.RotationRequest) (string, error) {
	return r.gen.Generate(req.GeneratorOpts)
}

func (r *Rotator) rotateKeyValue(ctx context.Context, req models.RotationRequest) (string, error) {
	existing, err := r.smClient.GetSecretValue(ctx, req.SecretARN)
	if err != nil {
		return "", fmt.Errorf("failed to get existing secret: %w", err)
	}

	var existingMap map[string]interface{}
	if err := json.Unmarshal([]byte(existing), &existingMap); err != nil {
		return "", fmt.Errorf("failed to parse existing secret as JSON: %w", err)
	}

	keysToRotate := getKeys(existingMap)
	if req.KeyValueConfig != nil && len(req.KeyValueConfig.KeysToRotate) > 0 {
		keysToRotate = req.KeyValueConfig.KeysToRotate
	}

	//newMap := make(map[string]interface{})
	for _, key := range keysToRotate {
		newValue, err := r.gen.Generate(req.GeneratorOpts)
		if err != nil {
			return "", fmt.Errorf("failed to generate value for key %s: %w", key, err)
		}
		existingMap[key] = newValue
	}

	result, err := json.Marshal(existingMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated secret: %w", err)
	}

	return string(result), nil
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
