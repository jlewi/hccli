package pkg

import (
	"context"
	"github.com/go-logr/zapr"
	"github.com/jlewi/hccli/pkg/config"
	"github.com/jlewi/hydros/pkg/files"
	"github.com/pkg/errors"
	"github.com/replicate/replicate-go"
	"go.uber.org/zap"
	"strings"
)

// ReplicateClient is a client for the model when deployed on Replicate (as opposed to K8s)
type ReplicateClient struct {
	client *replicate.Client
	config config.Config
}

func NewReplicateClient(cfg config.Config) (*ReplicateClient, error) {
	token, err := files.Read(cfg.Replicate.APITokenFile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read Replicate API token from file %s", cfg.Replicate.APITokenFile)
	}
	r8, err := replicate.NewClient(replicate.WithToken(strings.TrimSpace(string(token))))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create Replicate client")
	}
	return &ReplicateClient{
		client: r8,
		config: cfg,
	}, nil
}

// Translate takes a natural language query and returns the honeycomb query as a string
func (p *ReplicateClient) Translate(inQuery QueryInput) (string, error) {
	log := zapr.NewLogger(zap.L())

	input := replicate.PredictionInput{
		"nlq":  inQuery.NLQ,
		"cols": inQuery.COLS,
	}

	// Run a model and wait for its output
	log.Info("Sending prediction to Replicate", "query", inQuery)
	output, err := p.client.Run(context.Background(), p.config.Replicate.Model, input, nil)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to run model")
	}

	log.Info("Received query from Replicate", "output", output)

	outputQuery, ok := output.(string)

	if !ok {
		return "", errors.New("Failed to convert output to string")
	}

	return outputQuery, nil
}
