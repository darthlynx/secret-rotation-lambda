package rotator

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/darthlynx/secret-rotation-lambda/internal/models"
	"github.com/darthlynx/secret-rotation-lambda/internal/secretsmanager"
)

type mockGenerator struct {
	generateFunc func(opts models.GeneratorOptions) (string, error)
}

func (m *mockGenerator) Generate(opts models.GeneratorOptions) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(opts)
	}
	return "generated-secret", nil
}

func TestRotateSecret_Plaintext(t *testing.T) {
	mockSM := &secretsmanager.MockClient{
		PutSecretValueFunc: func(ctx context.Context, secretARN, secretValue string) (string, error) {
			return "version-123", nil
		},
	}

	mockGen := &mockGenerator{
		generateFunc: func(opts models.GeneratorOptions) (string, error) {
			return "new-plaintext-secret", nil
		},
	}

	rotator := New(mockSM, mockGen)

	req := models.RotationRequest{
		SecretARN:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:test",
		SecretType: models.SecretTypePlaintext,
		GeneratorOpts: models.GeneratorOptions{
			Length:           16,
			IncludeLowercase: true,
		},
	}

	resp, err := rotator.RotateSecret(context.Background(), req)

	if err != nil {
		t.Fatalf("RotateSecret() error: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success, got failure: %s", resp.ErrorMsg)
	}

	if resp.VersionID != "version-123" {
		t.Errorf("VersionID = %s, want version-123", resp.VersionID)
	}
}

func TestRotateSecret_KeyValue(t *testing.T) {
	existingSecret := map[string]interface{}{
		"username": "admin",
		"password": "old-password",
		"api_key":  "old-api-key",
	}
	existingJSON, _ := json.Marshal(existingSecret)

	var capturedValue string
	mockSM := &secretsmanager.MockClient{
		GetSecretValueFunc: func(ctx context.Context, secretARN string) (string, error) {
			return string(existingJSON), nil
		},
		PutSecretValueFunc: func(ctx context.Context, secretARN, secretValue string) (string, error) {
			capturedValue = secretValue
			return "version-456", nil
		},
	}

	callCount := 0
	mockGen := &mockGenerator{
		generateFunc: func(opts models.GeneratorOptions) (string, error) {
			callCount++
			return "new-generated-value", nil
		},
	}

	rotator := New(mockSM, mockGen)

	req := models.RotationRequest{
		SecretARN:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:test",
		SecretType: models.SecretTypeKeyValue,
		GeneratorOpts: models.GeneratorOptions{
			Length:           20,
			IncludeLowercase: true,
			IncludeDigits:    true,
		},
		KeyValueConfig: &models.KeyValueConfig{
			KeysToRotate: []string{"password", "api_key"},
		},
	}

	resp, err := rotator.RotateSecret(context.Background(), req)

	if err != nil {
		t.Fatalf("RotateSecret() error: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success, got failure: %s", resp.ErrorMsg)
	}

	if callCount != 2 {
		t.Errorf("Generator called %d times, want 2", callCount)
	}

	var updatedSecret map[string]interface{}
	if err := json.Unmarshal([]byte(capturedValue), &updatedSecret); err != nil {
		t.Fatalf("Failed to unmarshal updated secret: %v", err)
	}

	if updatedSecret["username"] != "admin" {
		t.Errorf("username should remain unchanged")
	}

	if updatedSecret["password"] != "new-generated-value" {
		t.Errorf("password should be updated")
	}
}

func TestRotateSecret_ValidationError(t *testing.T) {
	rotator := New(nil, nil)

	tests := []struct {
		name string
		req  models.RotationRequest
	}{
		{
			name: "missing secret ARN",
			req: models.RotationRequest{
				SecretType: models.SecretTypePlaintext,
			},
		},
		{
			name: "missing secret type",
			req: models.RotationRequest{
				SecretARN: "arn:aws:secretsmanager:us-east-1:123456789012:secret:test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := rotator.RotateSecret(context.Background(), tt.req)
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
			if resp.Success {
				t.Errorf("Expected failure response")
			}
		})
	}
}
