package client

import (
	"context"
	"encoding/json"
	"fmt"

	"net/url"

	"github.com/KhanhD1nh/langgraph-sdk-go/http"
	"github.com/KhanhD1nh/langgraph-sdk-go/schema"
)

type AssistantsClient struct {
	http *http.HttpClient
}

func NewAssistantsClient(httpClient *http.HttpClient) *AssistantsClient {
	return &AssistantsClient{http: httpClient}
}

func (c *AssistantsClient) Get(ctx context.Context, assistantID string, headers *map[string]string) (schema.Assistant, error) {
	result, err := c.http.Get(ctx, fmt.Sprintf("/assistants/%s", assistantID), nil, headers, nil)
	if err != nil {
		return schema.Assistant{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Assistant{}, err
	}

	var assistant schema.Assistant
	err = json.Unmarshal(jsonBytes, &assistant)
	if err != nil {
		return schema.Assistant{}, err
	}
	return assistant, nil
}

func (c *AssistantsClient) GetGraph(ctx context.Context, assistantID string, xray *bool, headers *map[string]string) (schema.Graph, error) {
	params := url.Values{}
	if xray != nil {
		params.Set("xray", fmt.Sprintf("%v", *xray))
	}

	result, err := c.http.Get(ctx, fmt.Sprintf("/assistants/%s/graph", assistantID), params, headers, nil)
	if err != nil {
		return schema.Graph{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Graph{}, err
	}

	var graph schema.Graph
	err = json.Unmarshal(jsonBytes, &graph)
	if err != nil {
		return schema.Graph{}, err
	}

	return graph, nil
}

func (c *AssistantsClient) GetSchemas(ctx context.Context, assistantID string, headers *map[string]string) (schema.GraphSchema, error) {
	result, err := c.http.Get(ctx, fmt.Sprintf("/assistants/%s/schemas", assistantID), nil, headers, nil)
	if err != nil {
		return schema.GraphSchema{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.GraphSchema{}, err
	}

	var graphSchema schema.GraphSchema
	err = json.Unmarshal(jsonBytes, &graphSchema)
	if err != nil {
		return schema.GraphSchema{}, err
	}

	return graphSchema, nil
}

func (c *AssistantsClient) GetSubgraphs(ctx context.Context, assistantID string, namespace *string, recurse *bool, headers *map[string]string) (schema.Subgraphs, error) {
	params := url.Values{}
	params.Set("recurse", fmt.Sprintf("%v", *recurse))

	var result interface{}
	var err error

	if namespace != nil {
		result, err = c.http.Get(ctx, fmt.Sprintf("/assistants/%s/subgraphs/%s", assistantID, *namespace), params, headers, nil)
		if err != nil {
			return schema.Subgraphs{}, err
		}
	} else {
		result, err = c.http.Get(ctx, fmt.Sprintf("/assistants/%s/subgraphs", assistantID), params, headers, nil)
		if err != nil {
			return schema.Subgraphs{}, err
		}
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Subgraphs{}, err
	}

	var subgraphs schema.Subgraphs
	err = json.Unmarshal(jsonBytes, &subgraphs)
	if err != nil {
		return schema.Subgraphs{}, err
	}

	return subgraphs, nil
}

func (c *AssistantsClient) Create(ctx context.Context, graphID *string, config *schema.Config, metadata *schema.Json, assistantID *string, ifExists *schema.OnConflictBehavior, name *string, headers *map[string]string, description *string) (schema.Assistant, error) {
	payload := map[string]any{
		"graph_id": graphID,
	}
	if config != nil {
		payload["config"] = *config
	}
	if metadata != nil {
		payload["metadata"] = *metadata
	}
	if assistantID != nil {
		payload["assistant_id"] = *assistantID
	}
	if ifExists != nil {
		payload["if_exists"] = *ifExists
	}
	if name != nil {
		payload["name"] = *name
	}
	if description != nil {
		payload["description"] = *description
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Post(ctx, "/assistants", payload, headers, nil)
	if err != nil {
		return schema.Assistant{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Assistant{}, err
	}

	var assistant schema.Assistant
	err = json.Unmarshal(jsonBytes, &assistant)
	if err != nil {
		return schema.Assistant{}, err
	}

	return assistant, nil
}

func (c *AssistantsClient) Update(ctx context.Context, assistantID string, graphID *string, config *schema.Config, metadata *schema.Json, name *string, headers *map[string]string, description *string) (schema.Assistant, error) {
	payload := map[string]any{}
	if graphID != nil {
		payload["graph_id"] = *graphID
	}
	if config != nil {
		payload["config"] = *config
	}
	if metadata != nil {
		payload["metadata"] = *metadata
	}
	if name != nil {
		payload["name"] = *name
	}
	if description != nil {
		payload["description"] = *description
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Patch(ctx, fmt.Sprintf("/assistants/%s", assistantID), payload, headers, nil)
	if err != nil {
		return schema.Assistant{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Assistant{}, err
	}

	var assistant schema.Assistant
	err = json.Unmarshal(jsonBytes, &assistant)
	if err != nil {
		return schema.Assistant{}, err
	}

	return assistant, nil
}

func (c *AssistantsClient) Delete(ctx context.Context, assistantID string, headers *map[string]string) error {
	err := c.http.Delete(ctx, fmt.Sprintf("/assistants/%s", assistantID), nil, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *AssistantsClient) Search(ctx context.Context, metadata *schema.Json, graphID *string, limit *int, offset *int, sortBy *schema.AssistantSortBy, sortOrder *schema.SortOrder, headers *map[string]string) ([]schema.Assistant, error) {
	if limit != nil && *limit <= 0 {
		*limit = 10
	}

	if offset != nil && *offset < 0 {
		*offset = 0
	}

	payload := map[string]any{
		"limit":  limit,
		"offset": offset,
	}
	if metadata != nil {
		payload["metadata"] = *metadata
	}
	if graphID != nil {
		payload["graph_id"] = *graphID
	}
	if sortBy != nil {
		payload["sort_by"] = *sortBy
	}
	if sortOrder != nil {
		payload["sort_order"] = *sortOrder
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Post(ctx, "/assistants/search", payload, headers, nil)
	if err != nil {
		return []schema.Assistant{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return []schema.Assistant{}, err
	}

	var assistants []schema.Assistant
	err = json.Unmarshal(jsonBytes, &assistants)
	if err != nil {
		return []schema.Assistant{}, err
	}

	return assistants, nil
}

func (c *AssistantsClient) GetVersions(ctx context.Context, assistantID string, metadata *schema.Json, limit *int, offset *int, headers *map[string]string) ([]schema.Assistant, error) {
	if limit != nil && *limit <= 0 {
		*limit = 10
	}

	if offset != nil && *offset < 0 {
		*offset = 0
	}

	payload := map[string]any{
		"limit":  limit,
		"offset": offset,
	}
	if metadata != nil {
		payload["metadata"] = *metadata
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Post(ctx, fmt.Sprintf("/assistants/%s/versions", assistantID), payload, headers, nil)
	if err != nil {
		return []schema.Assistant{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return []schema.Assistant{}, err
	}

	var assistants []schema.Assistant
	err = json.Unmarshal(jsonBytes, &assistants)
	if err != nil {
		return []schema.Assistant{}, err
	}

	return assistants, nil
}

func (c *AssistantsClient) SetLatest(ctx context.Context, assistantID string, version *int, headers *map[string]string) (schema.Assistant, error) {
	payload := map[string]any{
		"version": *version,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Post(ctx, fmt.Sprintf("/assistants/%s/versions/latest", assistantID), payload, headers, nil)
	if err != nil {
		return schema.Assistant{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Assistant{}, err
	}

	var assistant schema.Assistant
	err = json.Unmarshal(jsonBytes, &assistant)
	if err != nil {
		return schema.Assistant{}, err
	}

	return assistant, nil
}
