package http

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// decodeJSON decodes the response body as JSON
func decodeJSON(resp *resty.Response) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return result, nil
}
