package pkg

import (
	"encoding/json"
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
