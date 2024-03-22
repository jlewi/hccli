package pkg

import (
	"encoding/json"
	"github.com/jlewi/hccli/pkg/app"
	"os"
	"path/filepath"
	"testing"
)

const (
	datasetslug = "glider"
)

func Test_GetColumns(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skipf("Test is skipped in GitHub actions")
	}

	a := app.NewApp()
	if err := a.LoadConfig(nil); err != nil {
		t.Fatalf("Error loading config; %v", err)
	}

	if err := a.SetupLogging(); err != nil {
		t.Fatalf("Error setting up logging; %v", err)
	}

	hc, err := NewHoneycombClient(*a.Config)
	if err != nil {
		t.Fatalf("Error creating Honeycomb client; %v", err)
	}

	cols, err := hc.GetColumns(datasetslug)
	if err != nil {
		t.Fatalf("Error getting columns; %v", err)
	}
	t.Logf("Columns: %v", cols)
}

func Test_CreateQuery(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skipf("Test is skipped in GitHub actions")
	}

	a := app.NewApp()
	if err := a.LoadConfig(nil); err != nil {
		t.Fatalf("Error loading config; %v", err)
	}

	if err := a.SetupLogging(); err != nil {
		t.Fatalf("Error setting up logging; %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current directory; %v", err)
	}

	testQueryFile := filepath.Join(cwd, "test_data", "total_traces_query.json")
	queryb, err := os.ReadFile(testQueryFile)
	if err != nil {
		t.Fatalf("Error reading query file; %v", err)
	}
	query := &HoneycombQuery{}
	if err := json.Unmarshal(queryb, query); err != nil {
		t.Fatalf("Error unmarshalling query; %v", err)
	}

	hc, err := NewHoneycombClient(*a.Config)
	if err != nil {
		t.Fatalf("Error creating Honeycomb client; %v", err)
	}

	queryId, err := hc.CreateQuery(datasetslug, *query)
	if err != nil {
		t.Fatalf("Error creating query; %v", err)
	}
	t.Logf("Created query %s", queryId)
}
