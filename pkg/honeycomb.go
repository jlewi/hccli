package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/jlewi/hccli/pkg/config"
	"github.com/jlewi/hydros/pkg/files"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

const (
	honeycombEndpoint     = "https://api.honeycomb.ai"
	honeycombAPIKeyHeader = "X-Honeycomb-Team"
)

type HoneycombClient struct {
	apiKey string
}

func NewHoneycombClient(config config.Config) (*HoneycombClient, error) {
	apiKeyBytes, err := files.Read(config.HoneycombAPIKeyFile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read honeycomb API key from file %s", config.HoneycombAPIKeyFile)
	}

	apiKey := strings.TrimSpace(string(apiKeyBytes))
	if apiKey == "" {
		return nil, errors.New("Honeycomb API key is empty")
	}

	return &HoneycombClient{
		apiKey: apiKey,
	}, nil
}

func (h *HoneycombClient) GetColumns(datasetSlug string) (string, error) {
	log := zapr.NewLogger(zap.L())
	endpoint := fmt.Sprintf("https://api.honeycomb.io/1/columns/%s", datasetSlug)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create request")
	}

	req.Header.Set(honeycombAPIKeyHeader, h.apiKey)
	client := &http.Client{}

	log.Info("Fetching columns", "endpoint", endpoint)
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to send request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "Failed to read response body", "status", resp.StatusCode)
		} else {
			log.Info("Request failed", "status", resp.StatusCode, "body", string(body))

		}
		return "", errors.Errorf("Request failed with status code %v; body %v", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(body), errors.Wrapf(err, "Failed to read response body")
	}
	return string(body), nil
}

type HoneycombQuery struct {
	ID           *string  `json:"id,omitempty"`
	Breakdowns   []string `json:"breakdowns,omitempty"`
	Calculations []struct {
		Op     string      `json:"op,omitempty"`
		Column interface{} `json:"column,omitempty"`
	} `json:"calculations,omitempty"`
	Filters []struct {
		Op     string      `json:"op,omitempty"`
		Column interface{} `json:"column,omitempty"`
		Value  *struct {
		} `json:"value,omitempty"`
	} `json:"filters"`
	FilterCombination string `json:"filter_combination,omitempty"`
	Granularity       int    `json:"granularity,omitempty"`
	Orders            []struct {
		Column string `json:"column,omitempty"`
		Op     string `json:"op,omitempty"`
		Order  string `json:"order,omitempty"`
	} `json:"orders,omitempty"`
	Limit     int `json:"limit,omitempty"`
	StartTime int `json:"start_time,omitempty"`
	EndTime   int `json:"end_time,omitempty"`
	TimeRange int `json:"time_range,omitempty"`
	Havings   []struct {
		CalculateOp string      `json:"calculate_op,omitempty"`
		Column      interface{} `json:"column,omitempty"`
		Op          string      `json:"op,omitempty"`
		Value       int         `json:"value,omitempty"`
	} `json:"havings,omitempty"`
}

func (h *HoneycombClient) CreateQuery(datasetSlug string, q HoneycombQuery) (string, error) {
	log := zapr.NewLogger(zap.L())
	endpoint := fmt.Sprintf("https://api.honeycomb.io/1/queries/%s", datasetSlug)

	b, err := json.Marshal(q)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to serialize query")
	}
	buff := bytes.NewBuffer(b)
	req, err := http.NewRequest(http.MethodPost, endpoint, buff)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create request")
	}

	req.Header.Set(honeycombAPIKeyHeader, h.apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	log.Info("Creating query", "endpoint", endpoint, "query", string(b))
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to send request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "Failed to read response body", "status", resp.StatusCode)
		} else {
			log.Info("Request failed", "status", resp.StatusCode, "body", string(body))

		}
		return "", errors.Errorf("Request failed with status code %v; body %v", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to read response body")
	}

	outQuery := &HoneycombQuery{}
	if err := json.Unmarshal(body, outQuery); err != nil {
		return "", errors.Wrapf(err, "Failed to deserialize response body")
	}
	id := ""
	if outQuery.ID != nil {
		id = *outQuery.ID
	}
	return id, nil
}
