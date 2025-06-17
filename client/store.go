package client

import (
	"context"
	"encoding/json"
	"fmt"

	"net/url"
	"strings"

	"github.com/KhanhD1nh/langgraph-sdk-go/http"
	"github.com/KhanhD1nh/langgraph-sdk-go/schema"
)

type StoreClient struct {
	http *http.HttpClient
}

func NewStoreClient(httpClient *http.HttpClient) *StoreClient {
	return &StoreClient{http: httpClient}
}

func (c *StoreClient) PutItem(ctx context.Context, namespace []string, key string, value map[string]any, index *any, ttl *int, headers *map[string]string) error {
	for _, label := range namespace {
		if containsDot(label) {
			return fmt.Errorf("invalid namespace label '%s'. Namespace labels cannot contain periods ('.')", label)
		}
	}

	payload := map[string]any{
		"namespace": namespace,
		"key":       key,
		"value":     value,
		"index":     index,
		"ttl":       ttl,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	_, err := c.http.Put(ctx, "/store/items", payload, headers, nil)
	return err
}

func (c *StoreClient) GetItem(ctx context.Context, namespace []string, key string, refreshTtl *bool, headers *map[string]string) (map[string]any, error) {
	for _, label := range namespace {
		if containsDot(label) {
			return nil, fmt.Errorf("invalid namespace label '%s'. Namespace labels cannot contain periods ('.')", label)
		}
	}

	params := url.Values{}
	params.Add("namespace", strings.Join(namespace, "."))
	params.Add("key", key)

	if refreshTtl != nil {
		params.Add("refresh_ttl", fmt.Sprintf("%t", *refreshTtl))
	}

	result, err := c.http.Get(ctx, "/store/items", params, headers, nil)
	if err != nil {
		return nil, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var item map[string]any
	err = json.Unmarshal(jsonBytes, &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (c *StoreClient) DeleteItem(ctx context.Context, namespace []string, key string, headers *map[string]string) error {
	for _, label := range namespace {
		if containsDot(label) {
			return fmt.Errorf("invalid namespace label '%s'. Namespace labels cannot contain periods ('.')", label)
		}
	}

	jsonData := map[string]any{
		"namespace": namespace,
		"key":       key,
	}

	err := c.http.Delete(ctx, "/store/items", jsonData, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *StoreClient) SearchItems(ctx context.Context, namespacePrefix []string, filter *map[string]any, limit *int, offset *int, query *string, refreshTtl *bool, headers *map[string]string) (schema.SearchItemsResponse, error) {
	if limit != nil && *limit <= 0 {
		*limit = 10
	}

	if offset != nil && *offset < 0 {
		*offset = 0
	}

	payload := map[string]any{
		"namespace_prefix": namespacePrefix,
		"filter":           filter,
		"limit":            limit,
		"offset":           offset,
		"query":            query,
		"refresh_ttl":      refreshTtl,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Post(ctx, "/store/items/search", payload, headers, nil)
	if err != nil {
		return schema.SearchItemsResponse{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.SearchItemsResponse{}, err
	}

	var searchItemsResponse schema.SearchItemsResponse
	err = json.Unmarshal(jsonBytes, &searchItemsResponse)
	if err != nil {
		return schema.SearchItemsResponse{}, err
	}

	return searchItemsResponse, nil
}

func (c *StoreClient) ListNamespaces(ctx context.Context, prefix *[]string, suffix *[]string, maxDepth *int, limit *int, offset *int, headers *map[string]string) ([]schema.ListNamespaceResponse, error) {
	if limit != nil && *limit <= 0 {
		*limit = 100
	}

	if offset != nil && *offset < 0 {
		*offset = 0
	}

	payload := map[string]any{
		"prefix":    prefix,
		"suffix":    suffix,
		"max_depth": maxDepth,
		"limit":     limit,
		"offset":    offset,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Post(ctx, "/store/namespaces", payload, headers, nil)
	if err != nil {
		return []schema.ListNamespaceResponse{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return []schema.ListNamespaceResponse{}, err
	}

	var namespaces []schema.ListNamespaceResponse
	err = json.Unmarshal(jsonBytes, &namespaces)
	if err != nil {
		return []schema.ListNamespaceResponse{}, err
	}

	return namespaces, nil
}
