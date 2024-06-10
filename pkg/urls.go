package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/pkg/errors"
)

// QueryToURL converts a query to a URL.
// It uses Honeycomb's query parameter and template links feature
// https://docs.honeycomb.io/investigate/collaborate/share-query/
func QueryToURL(query HoneycombQuery, baseURL string, dataset string) (string, error) {
	b, err := json.Marshal(query)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to serialize query to JSON")
	}

	return baseURL + "/datasets/" + dataset + "?query=" + string(b), nil
}

// SaveHoneycombGraph saves a screenshot of a Honeycomb graph to a file
// It uses Chrome's remote debugging protocol to take a screenshot.
// You need to manually start a Chrome instance with remote debugging enabled
// e.g. chrome --remote-debugging-port=9222
// Then you must login into honeycomb.
// This will use your most recent chrome session so if you are logged into chrome with multiple accounts
// make sure the most recently used chrome is the one with the Honeycomb account.
func SaveHoneycombGraph(url string, outFile string, port int) error {
	log := zapr.NewLogger(zap.L())
	// Set up RemoteAllocator
	// You need to manually start a Chrome instance with remote debugging enabled
	// e.g. chrome --remote-debugging-port=9222
	// Then you can login.
	allocatorContext, cancelAllocator := chromedp.NewRemoteAllocator(context.Background(), fmt.Sprintf("http://localhost:%d", port))
	defer cancelAllocator()

	// Create context
	// If you don't call cancel the context stays open.
	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	// Run task list
	var buf []byte

	// TODO(jeremy): Will this selector work for other queries?
	svgSelector := `svg[data-testid="COUNT"]`
	divSelector := `div[data-testid="interaction-monitor"]`

	quality := 90

	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return errors.Wrapf(err, "Failed to navigate to the page and wait for it to be ready")
	}

	// N.B. We don't want to be pass a timeout to the first Run call because on the first Run call a browser is
	// created and the timeout would end up cancelling the browser and then the subsequent Run calls would fail.
	if err := chromedp.Run(ctx, waitForEither(svgSelector, divSelector, 7*time.Second)); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			// If there was a timeout then the selector is probably wrong so try to snapshot it any way.
			log.Info("Timed out waiting for page to be ready", "url", url, "divSelector", divSelector, "svgSelector", svgSelector)
		} else {
			return errors.Wrapf(err, "Failed to navigate to the page and wait for it to be ready")
		}
	}

	// Try to take the snapshot even if we got a timeout for the selector.
	if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, quality)); err != nil {
		return errors.Wrapf(err, "Failed to grab snapshot using chromedp")
	}
	// Save the screenshot to a file
	if err := os.WriteFile(outFile, buf, 0644); err != nil {
		return errors.Wrapf(err, "Failed to save screenshot to %s", outFile)
	}

	log.Info("Screenshot saved", "file", outFile)
	return nil
}

// waitForEither waits for either one of the selectors to be visible
func waitForEither(selector1, selector2 string, timeout time.Duration) chromedp.Tasks {
	// TODO(jeremy): This is chatGPT generated code. Does it make sense to wrap it in ChromeDPTasks?
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			defer cancel()

			done := make(chan error, 2)
			go func() {
				done <- chromedp.Run(ctx, chromedp.WaitReady(selector1, chromedp.ByQuery))
			}()
			go func() {
				done <- chromedp.Run(ctx, chromedp.WaitReady(selector2, chromedp.ByQuery))
			}()

			select {
			case err := <-done:
				cancel() // Cancel the other goroutine
				return err
			case <-time.After(timeout):
				cancel() // Cancel the other goroutine
				return context.DeadlineExceeded
			case <-ctx.Done():
				return ctx.Err()
			}
		}),
	}
}
