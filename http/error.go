package http

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func handleError(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.IsError() {
		body := resp.String()
		return fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), body)
	}
	return nil
}
