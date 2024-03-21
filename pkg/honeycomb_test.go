package pkg

import (
	"github.com/jlewi/hccli/pkg/app"
	"os"
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
