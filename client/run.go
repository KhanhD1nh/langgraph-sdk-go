package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/KhanhD1nh/langgraph-sdk-go/http"
	"github.com/KhanhD1nh/langgraph-sdk-go/schema"
)

type RunsClient struct {
	http *http.HttpClient
}

func NewRunsClient(httpClient *http.HttpClient) *RunsClient {
	return &RunsClient{http: httpClient}
}

func (c *RunsClient) Stream(ctx context.Context, threadID string, assistantID string, input *map[string]any, command *schema.Command, streamMode *[]schema.StreamMode, streamSubgraphs *bool, metadata *map[string]any, config *schema.Config, checkpoint *schema.Checkpoint, checkpointID *string, checkpointDuring *bool, interruptBefore *[]string, interruptAfter *[]string, feedbackKeys *[]string, webhook *string, multitaskStrategy *schema.MultitaskStrategy, ifNotExists *schema.IfNotExists, onDisconnect *schema.DisconnectMode, onCompletion *schema.OnCompletionBehavior, afterSeconds *int, headers *map[string]string) (chan schema.StreamPart, context.CancelFunc) {
	if streamMode != nil {
		streamMode = new([]schema.StreamMode)
		*streamMode = []schema.StreamMode{schema.StreamModeValues}
	}

	if streamSubgraphs != nil {
		streamSubgraphs = new(bool)
		*streamSubgraphs = false
	}

	payload := map[string]any{
		"input":              input,
		"command":            command,
		"config":             config,
		"metadata":           metadata,
		"stream_mode":        streamMode,
		"stream_subgraphs":   streamSubgraphs,
		"assistant_id":       assistantID,
		"interrupt_before":   interruptBefore,
		"interrupt_after":    interruptAfter,
		"feedback_keys":      feedbackKeys,
		"webhook":            webhook,
		"checkpoint":         checkpoint,
		"checkpoint_id":      checkpointID,
		"checkpoint_during":  checkpointDuring,
		"multitask_strategy": multitaskStrategy,
		"if_not_exists":      ifNotExists,
		"on_disconnect":      onDisconnect,
		"on_completion":      onCompletion,
		"after_seconds":      afterSeconds,
	}
	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	var endPoint string
	if threadID != "" {
		endPoint = fmt.Sprintf("/threads/%s/runs/stream", threadID)
	} else {
		endPoint = "/runs/stream"
	}

	ctx, cancel := context.WithCancel(ctx)

	streamCh, errCh, err := c.http.Stream(ctx, endPoint, "POST", payload, nil, headers, nil)
	if err != nil {
		cancel()
		errCh <- err
		close(errCh)
		return nil, cancel
	}

	return streamCh, cancel
}

func (c *RunsClient) Create(ctx context.Context, threadID string, assistantID string, input *map[string]any, command *schema.Command, streamMode *[]schema.StreamMode, streamSubgraphs *bool, metadata *map[string]any, config *schema.Config, checkpoint *schema.Checkpoint, checkpointID *string, checkpointDuring *bool, interruptBefore *[]string, interruptAfter *[]string, webhook *string, multitaskStrategy *schema.MultitaskStrategy, ifNotExists *schema.IfNotExists, onCompletion *schema.OnCompletionBehavior, afterSeconds *int, headers *map[string]string) (schema.Run, error) {
	payload := map[string]any{
		"input":              input,
		"command":            command,
		"config":             config,
		"metadata":           metadata,
		"stream_mode":        streamMode,
		"stream_subgraphs":   streamSubgraphs,
		"assistant_id":       assistantID,
		"interrupt_before":   interruptBefore,
		"interrupt_after":    interruptAfter,
		"webhook":            webhook,
		"checkpoint":         checkpoint,
		"checkpoint_id":      checkpointID,
		"checkpoint_during":  checkpointDuring,
		"multitask_strategy": multitaskStrategy,
		"if_not_exists":      ifNotExists,
		"after_seconds":      afterSeconds,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	var endPoint string
	if threadID != "" {
		endPoint = fmt.Sprintf("/threads/%s/runs", threadID)
	} else {
		endPoint = "/runs"
	}

	result, err := c.http.Post(ctx, endPoint, payload, headers, nil)
	if err != nil {
		return schema.Run{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Run{}, err
	}

	var run schema.Run
	err = json.Unmarshal(jsonBytes, &run)
	if err != nil {
		return schema.Run{}, err
	}

	return run, nil
}

func filterPayload(payload map[string]any) map[string]any {
	filtered := make(map[string]any)
	for k, v := range payload {
		if v != nil {
			filtered[k] = v
		}
	}
	return filtered
}

func (c *RunsClient) CreateBatch(ctx context.Context, payloads []map[string]any) ([]schema.Run, error) {
	filteredPayloads := make([]map[string]any, 0, len(payloads))
	for _, payload := range payloads {
		filteredPayloads = append(filteredPayloads, filterPayload(payload))
	}

	jsonData := map[string]any{"batch": filteredPayloads}

	result, err := c.http.Post(ctx, "/runs/batch", jsonData, nil, nil)
	if err != nil {
		return nil, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var runs []schema.Run
	err = json.Unmarshal(jsonBytes, &runs)
	if err != nil {
		return nil, err
	}

	return runs, nil
}

func (c *RunsClient) Wait(ctx context.Context, threadID string, assistantID string, input *map[string]any, command *schema.Command, metadata *map[string]any, config *schema.Config, checkPoint schema.Checkpoint, checkPointID *string, checkpointDuring *bool, interruptBefore *[]string, interruptAfter *[]string, webhook *string, onDisconnect *schema.DisconnectMode, onCompletion *schema.OnCompletionBehavior, multitaskStrategy *schema.MultitaskStrategy, ifNotExists *schema.IfNotExists, afterSeconds *int, raiseError *bool, headers *map[string]string) (any, error) {
	payload := map[string]any{
		"input":              input,
		"command":            command,
		"config":             config,
		"metadata":           metadata,
		"assistant_id":       assistantID,
		"checkpoint":         checkPoint,
		"checkpoint_id":      checkPointID,
		"checkpoint_during":  checkpointDuring,
		"interrupt_before":   interruptBefore,
		"interrupt_after":    interruptAfter,
		"webhook":            webhook,
		"multitask_strategy": multitaskStrategy,
		"if_not_exists":      ifNotExists,
		"on_disconnect":      onDisconnect,
		"on_completion":      onCompletion,
		"after_seconds":      afterSeconds,
		"raise_error":        raiseError,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	var endPoint string
	if threadID != "" {
		endPoint = fmt.Sprintf("/threads/%s/runs/wait", threadID)
	} else {
		endPoint = "/runs/wait"
	}

	result, err := c.http.Post(ctx, endPoint, payload, headers, nil)
	if err != nil {
		return nil, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var responseData map[string]any
	if err = json.Unmarshal(jsonBytes, &responseData); err != nil {
		return nil, err
	}

	if *raiseError {
		if errData, exists := responseData["__error__"].(map[string]any); exists {
			return nil, fmt.Errorf("%s", errData["message"])
		}
	}

	return responseData, nil
}

func (c *RunsClient) List(ctx context.Context, threadID string, limit *int, offset *int, status *schema.RunStatus, headers *map[string]string) ([]schema.Run, error) {
	if limit != nil && *limit <= 0 {
		*limit = 10
	}

	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("offset", fmt.Sprintf("%d", offset))

	if status != nil {
		params.Add("status", string(*status))
	}

	result, err := c.http.Get(ctx, fmt.Sprintf("/threads/%s/runs", threadID), params, headers, nil)
	if err != nil {
		return []schema.Run{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return []schema.Run{}, err
	}

	var runs []schema.Run
	err = json.Unmarshal(jsonBytes, &runs)
	if err != nil {
		return []schema.Run{}, err
	}

	return runs, nil
}

func (c *RunsClient) Get(ctx context.Context, threadID string, runID string, headers *map[string]string) (schema.Run, error) {
	result, err := c.http.Get(ctx, fmt.Sprintf("/threads/%s/runs/%s", threadID, runID), nil, headers, nil)
	if err != nil {
		return schema.Run{}, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return schema.Run{}, err
	}

	var run schema.Run
	err = json.Unmarshal(jsonBytes, &run)
	if err != nil {
		return schema.Run{}, err
	}

	return run, nil
}

func (c *RunsClient) Cancel(ctx context.Context, threadID string, runID string, wait *bool, action *schema.CancelAction, headers *map[string]string) error {
	if action != nil && *action == "" {
		*action = schema.CancelActionInterrupt
	}

	payload := map[string]any{
		"wait":   wait,
		"action": action,
	}

	payload, ok := removeEmptyFields(payload).(map[string]any)
	if !ok {
		fmt.Println("Error: cleanedPayload is not a map[string]any")
	}

	_, err := c.http.Post(ctx, fmt.Sprintf("/threads/%s/runs/%s/cancel", threadID, runID), payload, headers, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *RunsClient) Join(ctx context.Context, threadID string, runID string, headers *map[string]string) (map[string]any, error) {
	result, err := c.http.Get(ctx, fmt.Sprintf("/threads/%s/runs/%s/join", threadID, runID), nil, headers, nil)
	if err != nil {
		return nil, err
	}

	// Convert result to JSON bytes
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var responseData map[string]any
	err = json.Unmarshal(jsonBytes, &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func (c *RunsClient) JoinStream(ctx context.Context, threadID string, runID string, cancelOnDisconnect *bool, streamMode *[]schema.StreamMode, headers *map[string]string) (chan schema.StreamPart, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	params := url.Values{}

	if cancelOnDisconnect != nil {
		cancelOnDisconnect = new(bool)
		*cancelOnDisconnect = false
		params.Add("cancel_on_disconnect", fmt.Sprintf("%t", *cancelOnDisconnect))
	}
	if streamMode != nil {
		streamMode = new([]schema.StreamMode)
		*streamMode = []schema.StreamMode{schema.StreamModeValues}
		params.Add("stream_mode", fmt.Sprintf("%s", *streamMode))
	}

	streamCh, errCh, err := c.http.Stream(ctx, fmt.Sprintf("/threads/%s/runs/%s/join/stream", threadID, runID), "GET", nil, params, headers, nil)
	if err != nil {
		cancel()
		errCh <- err
		close(errCh)
		return nil, cancel
	}

	convertedCh := make(chan schema.StreamPart)

	go func() {
		defer close(convertedCh)
		for streamPart := range streamCh {
			convertedCh <- schema.StreamPart{
				Event: streamPart.Event,
				Data:  streamPart.Data,
			}
		}
	}()

	return convertedCh, cancel
}

func (c *RunsClient) Delete(ctx context.Context, threadID string, runID string, headers *map[string]string) error {
	err := c.http.Delete(ctx, fmt.Sprintf("/threads/%s/runs/%s", threadID, runID), nil, headers, nil)
	if err != nil {
		return err
	}

	return nil
}
