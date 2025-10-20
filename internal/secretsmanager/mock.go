package secretsmanager

import (
	"context"
	"errors"
)

// MockClient is a mock implementation of the Client interface for testing.
type MockClient struct {
	GetSecretValueFunc func(ctx context.Context, secretARN string) (string, error)
	PutSecretValueFunc func(ctx context.Context, secretARN, secretValue string) (string, error)
}

// GetSecretValue calls the mock function.
func (m *MockClient) GetSecretValue(ctx context.Context, secretARN string) (string, error) {
	if m.GetSecretValueFunc != nil {
		return m.GetSecretValueFunc(ctx, secretARN)
	}
	return "", errors.New("GetSecretValueFunc not implemented")
}

// PutSecretValue calls the mock function.
func (m *MockClient) PutSecretValue(ctx context.Context, secretARN, secretValue string) (string, error) {
	if m.PutSecretValueFunc != nil {
		return m.PutSecretValueFunc(ctx, secretARN, secretValue)
	}
	return "", errors.New("PutSecretValueFunc not implemented")
}
