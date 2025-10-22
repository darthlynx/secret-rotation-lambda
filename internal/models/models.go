package models

type SecretType string

const (
	SecretTypePlaintext SecretType = "plaintext"
	SecretTypeKeyValue  SecretType = "key-value"
	SecretTypeJSON      SecretType = "json"
)

// RotationRequest represents the input parameters for secret rotation.
type RotationRequest struct {
	SecretARN      string           `json:"secret_arn"`
	SecretType     SecretType       `json:"secret_type"`
	GeneratorOpts  GeneratorOptions `json:"generator_options"`
	KeyValueConfig *KeyValueConfig  `json:"key_value_config,omitempty"`
}

// GeneratorOptions defines options for secret generation.
type GeneratorOptions struct {
	Length              int  `json:"length"`
	IncludeDigits       bool `json:"include_digits"`
	IncludeUppercase    bool `json:"include_uppercase"`
	IncludeSpecialChars bool `json:"include_special_chars"`
}

// KeyValueConfig speicifies which keys to rotate in the key-value secrets
type KeyValueConfig struct {
	KeysToRotate []string `json:"keys_to_rotate"` // empty means rotate all
}

// RotationResponse represents the result of a secret rotation operation.
type RotationResponse struct {
	Success   bool   `json:"success"`
	SecretARN string `json:"secret_arn"`
	VersionID string `json:"version_id,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}
