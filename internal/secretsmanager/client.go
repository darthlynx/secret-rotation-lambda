package secretsmanager

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Client defines the interface for AWS Secrets Manager operations.
type Client interface {
	GetSecretValue(ctx context.Context, secretARN string) (string, error)
	PutSecretValue(ctx context.Context, secretARN, secretValue string) (string, error)
}

// SecretsManagerClient implements the Client interface.
type SecretsManagerClient struct {
	client *secretsmanager.Client
}

// NewClient creates a new SecretsManagerClient.
func NewClient(cfg aws.Config) *SecretsManagerClient {
	return &SecretsManagerClient{
		client: secretsmanager.NewFromConfig(cfg),
	}
}

// GetSecretValue retrieves the current secret value.
func (c *SecretsManagerClient) GetSecretValue(ctx context.Context, secretARN string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretARN),
	}

	result, err := c.client.GetSecretValue(ctx, input)
	if err != nil {
		return "", err
	}

	if result.SecretString != nil {
		return *result.SecretString, nil
	}

	return "", nil
}

// PutSecretValue updates the secret with a new value.
func (c *SecretsManagerClient) PutSecretValue(ctx context.Context, secretARN, secretValue string) (string, error) {
	input := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(secretARN),
		SecretString: aws.String(secretValue),
	}

	result, err := c.client.PutSecretValue(ctx, input)
	if err != nil {
		return "", err
	}

	return *result.VersionId, nil
}

// MergeKeyValueSecret merges new keys into existing key-value secret.
func MergeKeyValueSecret(existing, newValues string, keysToRotate []string) (string, error) {
	var existingMap, newMap map[string]interface{}

	if err := json.Unmarshal([]byte(existing), &existingMap); err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(newValues), &newMap); err != nil {
		return "", err
	}

	// If keysToRotate is empty, rotate all keys from newMap
	if len(keysToRotate) == 0 {
		for k, v := range newMap {
			existingMap[k] = v
		}
	} else {
		// Only rotate specified keys
		for _, key := range keysToRotate {
			if val, ok := newMap[key]; ok {
				existingMap[key] = val
			}
		}
	}

	merged, err := json.Marshal(existingMap)
	if err != nil {
		return "", err
	}

	return string(merged), nil
}
