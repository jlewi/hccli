package pkg

import (
	"encoding/json"
	"github.com/jlewi/hccli/pkg/config"
	"go.uber.org/zap"
	"os"
	"testing"
)

func Test_Replicate(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skipf("Test is skipped in GitHub actions")
	}

	if err := config.InitViper(nil); err != nil {
		t.Fatalf("Failed to initialize viper: %v", err)
	}
	cfg := config.GetConfig()

	log, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	zap.ReplaceGlobals(log)

	client, err := NewReplicateClient(*cfg)

	if err != nil {
		t.Fatalf("Failed to create Replicate client: %v", err)
	}

	cols := []string{"sli.latency", "duration_ms", "trace.trace_id", "trace.span_id", "trace.parent_id"}
	jsonCols, err := json.Marshal(cols)
	if err != nil {
		t.Fatalf("Failed to marshal columns: %v", err)
	}
	input := QueryInput{
		NLQ:  "Traces for the last 7 days",
		COLS: string(jsonCols),
	}
	query, err := client.Translate(input)
	if err != nil {
		t.Fatalf("Failed to predict: %v", err)
	}

	t.Logf("Query: %+v", query)
}
