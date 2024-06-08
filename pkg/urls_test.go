package pkg

import (
	"os"
	"testing"
)

func Test_ChromeDP(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skipf("Test is skipped in GitHub actions")
	}

	outFile := "screenshot.png"
	url := "https://ui.honeycomb.io/autobuilder/environments/prod/datasets/autobuilder/result/mm2wZinaKtT"
	port := 9222
	if err := SaveHoneycombGraph(url, outFile, port); err != nil {
		t.Fatalf("Failed to run; %v", err)
	}
}
