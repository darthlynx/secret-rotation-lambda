package secretsmanager

import (
	"encoding/json"
	"testing"
)

func TestMergeKeyValueSecret(t *testing.T) {
	tests := []struct {
		name         string
		existing     string
		newValues    string
		keysToRotate []string
		wantKeys     map[string]string
		wantErr      bool
	}{
		{
			name: "rotate all keys",
			existing: `{
				"username": "admin",
				"password": "old-pass",
				"api_key": "old-key"
			}`,
			newValues: `{
				"password": "new-pass",
				"api_key": "new-key"
			}`,
			keysToRotate: []string{},
			wantKeys: map[string]string{
				"username": "admin",
				"password": "new-pass",
				"api_key":  "new-key",
			},
			wantErr: false,
		},
		{
			name: "rotate specific keys",
			existing: `{
				"username": "admin",
				"password": "old-pass",
				"api_key": "old-key"
			}`,
			newValues: `{
				"password": "new-pass",
				"api_key": "new-key"
			}`,
			keysToRotate: []string{"password"},
			wantKeys: map[string]string{
				"username": "admin",
				"password": "new-pass",
				"api_key":  "old-key",
			},
			wantErr: false,
		},
		{
			name:         "invalid existing JSON",
			existing:     `invalid json`,
			newValues:    `{"key": "value"}`,
			keysToRotate: []string{},
			wantErr:      true,
		},
		{
			name:         "invalid new JSON",
			existing:     `{"key": "value"}`,
			newValues:    `invalid json`,
			keysToRotate: []string{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MergeKeyValueSecret(tt.existing, tt.newValues, tt.keysToRotate)

			if tt.wantErr {
				if err == nil {
					t.Errorf("MergeKeyValueSecret() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("MergeKeyValueSecret() unexpected error: %v", err)
				return
			}

			var resultMap map[string]interface{}
			if err := json.Unmarshal([]byte(result), &resultMap); err != nil {
				t.Fatalf("Failed to unmarshal result: %v", err)
			}

			for key, expectedValue := range tt.wantKeys {
				actualValue, ok := resultMap[key]
				if !ok {
					t.Errorf("Key %s not found in result", key)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("Key %s: got %v, want %v", key, actualValue, expectedValue)
				}
			}
		})
	}
}
