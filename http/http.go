package http

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/KhanhD1nh/langgraph-sdk-go/schema"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// ResponseCallback is a function that will be called with the response
type ResponseCallback func(*resty.Response)

// HttpClient handles async requests to the LangGraph API.
// Adds additional error messaging & content handling above the
// provided resty client.
type HttpClient struct {
	client *resty.Client
}

// NewHttpClient creates a new HttpClient with resty.Client
func NewHttpClient(baseURL string, headers map[string]string, timeOut time.Duration, transport http.RoundTripper) *HttpClient {
	client := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Accept", "application/json").
		SetHeaders(headers).
		SetTimeout(timeOut).
		SetTransport(transport)
	return &HttpClient{
		client: client,
	}
}

func (c *HttpClient) CheckConnection() error {
	_, err := c.client.R().Get("/")
	return err
}

// Get sends a GET request.
func (c *HttpClient) Get(ctx context.Context, path string, params url.Values, header *map[string]string, onResponse ResponseCallback) (interface{}, error) {
	req := c.client.R().SetContext(ctx)
	if params != nil {
		req.SetQueryParamsFromValues(params)
	}

	if header != nil {
		for key, value := range *header {
			req.SetHeader(key, value)
		}
	}
	resp, err := req.Get(path)

	if onResponse != nil {
		onResponse(resp)
	}

	if err := handleError(resp, err); err != nil {
		return nil, err
	}

	return decodeJSON(resp)
}

// Post sends a POST request.
func (c *HttpClient) Post(ctx context.Context, path string, jsonData any, header *map[string]string, onResponse ResponseCallback) (interface{}, error) {
	req := c.client.R().SetContext(ctx)

	if jsonData != nil {
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(jsonData)
	}

	if header != nil {
		for key, value := range *header {
			req.SetHeader(key, value)
		}
	}

	resp, err := req.Post(path)

	if onResponse != nil {
		onResponse(resp)
	}

	if err := handleError(resp, err); err != nil {
		return nil, err
	}

	return decodeJSON(resp)
}

// Put sends a PUT request.
func (c *HttpClient) Put(ctx context.Context, path string, jsonData any, header *map[string]string, onResponse ResponseCallback) (interface{}, error) {
	req := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData)

	if header != nil {
		for key, value := range *header {
			req.SetHeader(key, value)
		}
	}

	resp, err := req.Put(path)

	if onResponse != nil {
		onResponse(resp)
	}

	if err := handleError(resp, err); err != nil {
		return nil, err
	}

	return decodeJSON(resp)
}

// Patch sends a PATCH request.
func (c *HttpClient) Patch(ctx context.Context, path string, jsonData any, header *map[string]string, onResponse ResponseCallback) (interface{}, error) {
	req := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData)

	if header != nil {
		for key, value := range *header {
			req.SetHeader(key, value)
		}
	}

	resp, err := req.Patch(path)

	if onResponse != nil {
		onResponse(resp)
	}

	if err := handleError(resp, err); err != nil {
		return nil, err
	}

	return decodeJSON(resp)
}

// Delete sends a DELETE request.
func (c *HttpClient) Delete(ctx context.Context, path string, jsonData any, header *map[string]string, onResponse ResponseCallback) error {
	req := c.client.R().SetContext(ctx)

	if jsonData != nil {
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(jsonData)
	}

	if header != nil {
		for key, value := range *header {
			req.SetHeader(key, value)
		}
	}

	resp, err := req.Delete(path)

	if onResponse != nil {
		onResponse(resp)
	}

	if err := handleError(resp, err); err != nil {
		return err
	}

	return nil
}

// Stream streams results using SSE.
func (c *HttpClient) Stream(ctx context.Context, path string, method string, jsonData any, params url.Values, headers *map[string]string, onResponse ResponseCallback) (chan schema.StreamPart, chan error, error) {
	req := c.client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true). // Important for streaming
		SetHeader("Accept", "text/event-stream").
		SetHeader("Cache-Control", "no-store")

	if headers != nil {
		for key, value := range *headers {
			c.client.R().SetHeader(key, value)
		}
	}

	if jsonData != nil {
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(jsonData)
	}

	if params != nil {
		req.SetQueryParamsFromValues(params)
	}

	var resp *resty.Response
	var err error

	// Execute request based on method
	switch strings.ToUpper(method) {
	case "GET":
		resp, err = req.Get(path)
	case "POST":
		resp, err = req.Post(path)
	case "PUT":
		resp, err = req.Put(path)
	case "PATCH":
		resp, err = req.Patch(path)
	case "DELETE":
		resp, err = req.Delete(path)
	default:
		return nil, nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if onResponse != nil {
		onResponse(resp)
	}

	if err != nil {
		return nil, nil, err
	}

	// Get raw response body
	rawBody := resp.RawBody()

	// Check status code
	if resp.StatusCode() >= 400 {
		// Read error body
		body, _ := io.ReadAll(rawBody)
		rawBody.Close()
		return nil, nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), string(body))
	}

	// Check content type
	contentType := resp.Header().Get("Content-Type")
	if contentType == "" || !containsTextEventStream(contentType) {
		rawBody.Close()
		return nil, nil, fmt.Errorf("expected Content-Type to contain 'text/event-stream', got %s", contentType)
	}

	streamPartCh := make(chan schema.StreamPart)
	errCh := make(chan error, 1)

	// Process the SSE stream in a goroutine
	go func() {
		defer rawBody.Close()
		defer close(streamPartCh)
		defer close(errCh)

		// Parse SSE manually, since you mentioned seeing valid SSE data
		scanner := bufio.NewScanner(rawBody)
		var event, data, metadata string

		for scanner.Scan() {
			line := scanner.Text()
			// Empty line marks the end of an event
			if line == "" {
				if event != "" || data != "" {
					streamPartCh <- schema.StreamPart{
						Event:    event,
						Data:     data,
						MetaData: metadata,
					}
					// Reset for next event
					event = ""
					data = ""
					metadata = ""
				}
				continue
			} else {
				event = gjson.Get(line, "event").String()
				data = gjson.Get(line, "data").Raw
				metadata = gjson.Get(line, "metadata").Raw

				if event != "" || data != "" || metadata != "" {
					streamPartCh <- schema.StreamPart{
						Event:    event,
						Data:     data,
						MetaData: metadata,
					}
				}

				// Reset for next event
				event = ""
				data = ""
				metadata = ""
			}
		}
	}()

	return streamPartCh, errCh, nil
}
