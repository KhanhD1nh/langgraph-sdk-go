package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/KhanhD1nh/langgraph-sdk-go/http"
	"github.com/KhanhD1nh/langgraph-sdk-go/schema"
)

type ThreadsClient struct {
	http *http.HttpClient
}

func NewThreadsClient(httpClient *http.HttpClient) *ThreadsClient {
	return &ThreadsClient{http: httpClient}
}

func (c *ThreadsClient) Get(ctx context.Context, threadID string, headers *map[string]string) (schema.Thread, error) {
	result, err := c.http.Get(ctx, fmt.Sprintf("/threads/%s", threadID), nil, headers, nil)
	if err != nil {
		return schema.Thread{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Thread{}, err
	}

	var thread schema.Thread
	err = json.Unmarshal(jsonBytes, &thread)
	if err != nil {
		return schema.Thread{}, err
	}

	return thread, nil
}

func (c *ThreadsClient) Create(ctx context.Context, metadata *schema.Json, threadID *string, ifExists *schema.OnConflictBehavior, supersteps *[]any, graphID *string, headers *map[string]string) (schema.Thread, error) {
	payload := map[string]any{}
	if metadata != nil {
		payload["metadata"] = *metadata
	}
	if threadID != nil {
		payload["thread_id"] = *threadID
	}
	if ifExists != nil {
		payload["if_exists"] = *ifExists
	}
	if supersteps != nil {
		superstepsSlice := *supersteps
		var superstepsPayload []map[string]any
		for _, s := range superstepsSlice {
			sMap, ok := s.(map[string]any)
			if !ok {
				return schema.Thread{}, fmt.Errorf("each superstep must be a map, got %T", s)
			}
			updatesRaw, ok := sMap["updates"]
			if !ok {
				return schema.Thread{}, fmt.Errorf("superstep missing 'updates' key")
			}
			updatesSlice, ok := updatesRaw.([]any)
			if !ok {
				return schema.Thread{}, fmt.Errorf("'updates' must be a slice, got %T", updatesRaw)
			}
			var updatesPayload []map[string]any
			for _, u := range updatesSlice {
				uMap, ok := u.(map[string]any)
				if !ok {
					return schema.Thread{}, fmt.Errorf("each update must be a map, got %T", u)
				}
				updateObj := map[string]any{
					"values":  uMap["values"],
					"as_node": uMap["as_node"],
				}
				if cmd, ok := uMap["command"]; ok {
					updateObj["command"] = cmd
				}
				updatesPayload = append(updatesPayload, updateObj)
			}
			superstepsPayload = append(superstepsPayload, map[string]any{
				"updates": updatesPayload,
			})
		}
		payload["supersteps"] = superstepsPayload
	}
	if graphID != nil {
		payload["graph_id"] = *graphID
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Post(ctx, "/threads", payload, headers, nil)
	if err != nil {
		return schema.Thread{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Thread{}, err
	}

	var thread schema.Thread
	err = json.Unmarshal(jsonBytes, &thread)
	if err != nil {
		return schema.Thread{}, err
	}

	return thread, nil
}

func (c *ThreadsClient) Update(ctx context.Context, threadID string, metadata *schema.Json, headers *map[string]string) (schema.Thread, error) {
	payload := map[string]any{}
	if metadata != nil {
		payload["metadata"] = *metadata
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Patch(ctx, fmt.Sprintf("/threads/%s", threadID), payload, headers, nil)
	if err != nil {
		return schema.Thread{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Thread{}, err
	}

	var thread schema.Thread
	err = json.Unmarshal(jsonBytes, &thread)
	if err != nil {
		return schema.Thread{}, err
	}

	return thread, nil
}

func (c *ThreadsClient) Delete(ctx context.Context, threadID string, headers *map[string]string) error {
	err := c.http.Delete(ctx, fmt.Sprintf("/threads/%s", threadID), nil, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *ThreadsClient) Search(ctx context.Context, metadata *schema.Json, values *schema.Json, status *schema.ThreadStatus, limit *int, offset *int, sortBy *schema.ThreadSortBy, sortOrder *schema.SortOrder, headers *map[string]string) ([]schema.Thread, error) {
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
		payload["metadata"] = metadata
	}
	if values != nil {
		payload["values"] = values
	}
	if status != nil {
		payload["status"] = *status
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

	result, err := c.http.Post(ctx, "/threads/search", payload, headers, nil)
	if err != nil {
		return []schema.Thread{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return []schema.Thread{}, err
	}

	var threads []schema.Thread
	err = json.Unmarshal(jsonBytes, &threads)
	if err != nil {
		return []schema.Thread{}, err
	}

	return threads, nil
}

func (c *ThreadsClient) Copy(ctx context.Context, threadID string, headers *map[string]string) error {
	_, err := c.http.Post(ctx, fmt.Sprintf("/threads/%s/copy", threadID), nil, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *ThreadsClient) GetState(ctx context.Context, threadID string, checkPoint *schema.Checkpoint, checkPointID *string, subgraphs *bool, headers *map[string]string) (schema.ThreadState, error) {
	if subgraphs == nil {
		subgraphs = new(bool)
		*subgraphs = false
	}

	if checkPoint != nil {
		payload := map[string]any{
			"checkpoint": *checkPoint,
			"subgraphs":  *subgraphs,
		}

		payload, ok := removeEmptyFields(payload).(map[string]any)
		if !ok {
			fmt.Println("Error: cleanedPayload is not a map[string]any")
		}

		result, err := c.http.Post(ctx, fmt.Sprintf("/threads/%s/state/checkpoint", threadID), payload, headers, nil)
		if err != nil {
			return schema.ThreadState{}, err
		}

		// Convert result to JSON bytes
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return schema.ThreadState{}, err
		}

		var threadState schema.ThreadState
		err = json.Unmarshal(jsonBytes, &threadState)
		if err != nil {
			return schema.ThreadState{}, err
		}

		return threadState, nil
	} else if checkPointID != nil {
		result, err := c.http.Get(ctx, fmt.Sprintf("/threads/%s/state/%s", threadID, *checkPointID), nil, headers, nil)
		if err != nil {
			return schema.ThreadState{}, err
		}

		// Convert result to JSON bytes
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return schema.ThreadState{}, err
		}

		var threadState schema.ThreadState
		err = json.Unmarshal(jsonBytes, &threadState)
		if err != nil {
			return schema.ThreadState{}, err
		}

		return threadState, nil
	} else {
		result, err := c.http.Get(ctx, fmt.Sprintf("/threads/%s/state", threadID), nil, headers, nil)
		if err != nil {
			return schema.ThreadState{}, err
		}

		// Convert result to JSON bytes
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return schema.ThreadState{}, err
		}

		var threadState schema.ThreadState
		err = json.Unmarshal(jsonBytes, &threadState)
		if err != nil {
			return schema.ThreadState{}, err
		}

		return threadState, nil
	}
}

func (c *ThreadsClient) UpdateState(ctx context.Context, threadID string, values *any, asNode *string, checkPoint *schema.Checkpoint, checkPointID *string, headers *map[string]string) (schema.ThreadUpdateStateResponse, error) {
	payload := map[string]any{
		"values": *values,
	}
	if asNode != nil {
		payload["as_node"] = *asNode
	}
	if checkPoint != nil {
		payload["checkpoint"] = *checkPoint
	}
	if checkPointID != nil {
		payload["checkpoint_id"] = *checkPointID
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Post(ctx, fmt.Sprintf("/threads/%s/state", threadID), payload, headers, nil)
	if err != nil {
		return schema.ThreadUpdateStateResponse{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.ThreadUpdateStateResponse{}, err
	}

	var threadUpdateStateResponse schema.ThreadUpdateStateResponse
	err = json.Unmarshal(jsonBytes, &threadUpdateStateResponse)
	if err != nil {
		return schema.ThreadUpdateStateResponse{}, err
	}

	return threadUpdateStateResponse, nil
}

func (c *ThreadsClient) GetHistory(ctx context.Context, threadID string, limit *int, before *any, metadata *map[string]any, checkPoint *schema.Checkpoint, headers *map[string]string) ([]schema.ThreadState, error) {
	if limit != nil && *limit <= 0 {
		*limit = 10
	}

	payload := map[string]any{
		"limit": limit,
	}
	if before != nil {
		payload["before"] = before
	}
	if metadata != nil {
		payload["metadata"] = metadata
	}
	if checkPoint != nil {
		payload["checkpoint"] = *checkPoint
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	result, err := c.http.Post(ctx, fmt.Sprintf("/threads/%s/history", threadID), payload, headers, nil)
	if err != nil {
		return []schema.ThreadState{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return []schema.ThreadState{}, err
	}

	var threadStates []schema.ThreadState
	err = json.Unmarshal(jsonBytes, &threadStates)
	if err != nil {
		return []schema.ThreadState{}, err
	}

	return threadStates, nil
}
