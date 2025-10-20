package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/darthlynx/secret-rotation-lambda/internal/generator"
	"github.com/darthlynx/secret-rotation-lambda/internal/models"
	"github.com/darthlynx/secret-rotation-lambda/internal/rotator"
	"github.com/darthlynx/secret-rotation-lambda/internal/secretsmanager"
)

var rot *rotator.Rotator

func init() {
	logger := slog.New(slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("Failed to load AWS config", "err", err)
		return
	}
	smClient := secretsmanager.NewClient(cfg)
	gen := generator.New()
	rot = rotator.New(smClient, gen)
}

// HandleRequest is the Lambda function handler
func HandleRequest(ctx context.Context, event json.RawMessage) (*models.RotationResponse, error) {
	var req models.RotationRequest
	if err := json.Unmarshal(event, &req); err != nil {
		return &models.RotationResponse{
			Success:  false,
			ErrorMsg: "Invalid request format: " + err.Error(),
		}, err
	}

	return rot.RotateSecret(ctx, req)
}

func main() {
	lambda.Start(HandleRequest)
}
