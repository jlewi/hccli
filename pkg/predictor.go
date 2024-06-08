package pkg

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/zapr"
	"github.com/jlewi/hccli/pkg/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Translator translates a natural language query to a honeycomb query.
type Translator interface {
	Translate(nlq QueryInput) (string, error)
}

// Query represents a query for the model; not a honeycomb query.
type Query struct {
	Input               *QueryInput  `json:"input,omitempty"`
	Output              *string      `json:"output,omitempty"`
	Id                  *string      `json:"id,omitempty"`
	Version             *string      `json:"version,omitempty"`
	CreatedAt           *time.Time   `json:"created_at,omitempty"`
	StartedAt           *time.Time   `json:"started_at,omitempty"`
	CompletedAt         *time.Time   `json:"completed_at,omitempty"`
	Logs                *string      `json:"logs,omitempty"`
	Error               *interface{} `json:"error,omitempty"`
	Status              *string      `json:"status,omitempty"`
	Metrics             Metrics      `json:"metrics,omitempty"`
	OutputFilePrefix    *string      `json:"output_file_prefix,omitempty"`
	Webhook             *string      `json:"webhook,omitempty"`
	WebhookEventsFilter []string     `json:"webhook_events_filter,omitempty"`
}

type QueryInput struct {
	NLQ string `json:"nlq"`
	// cols is a string representing a json list of cols.
	COLS string `json:"cols"`
}

type Metrics struct {
	PredictTime float64 `json:"predict_time"`
}

type Predictor struct {
	Config *config.Config
}

func (p *Predictor) Predict(inQuery QueryInput) (*Query, error) {
	log := zapr.NewLogger(zap.L())
	q := &Query{
		Input: &inQuery,
	}
	b, err := json.Marshal(q)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to serialize query")
	}
	buff := bytes.NewBuffer(b)

	endpoint := strings.TrimSuffix(p.Config.AIEndpoint, "/") + "/predictions"
	log.Info("Sending prediction request", "endpoint", endpoint, "query", string(b))
	req, err := http.NewRequest(http.MethodPost, endpoint, buff)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to send request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "Failed to read response body", "status", resp.StatusCode)
		} else {
			log.Info("Request failed", "status", resp.StatusCode, "body", string(body))

		}
		return nil, errors.Errorf("Request failed with status code %v; body %v", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read response body")
	}
	outQuery := &Query{}

	if err := json.Unmarshal(body, outQuery); err != nil {
		return nil, errors.Wrapf(err, "Failed to deserialize response body")
	}
	return outQuery, nil
}

func (p *Predictor) Translate(inQuery QueryInput) (string, error) {
	query, err := p.Predict(inQuery)
	if err != nil {
		return "", err
	}
	return *query.Output, err
}
