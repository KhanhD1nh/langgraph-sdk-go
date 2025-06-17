// The LangGraph client implementations connect to the LangGraph API.
//
// This module provides both asynchronous (LangGraphClient) and synchronous (SyncLanggraphClient)
// clients to interacting with the LangGraph API's core resources such as
// Assistants, Threads, Runs, and Cron jobs, as well as its persistent
// document Store.

package langgraph_sdk

import (
	"fmt"

	"maps"
	http_client "net/http"
	"os"
	"strings"
	"time"

	"github.com/KhanhD1nh/langgraph-sdk-go/client"
	"github.com/KhanhD1nh/langgraph-sdk-go/http"
)

var (
	RESERVED_HEADERS = []string{"x-api-key"}
	Version          = "v1.0.5"
)

type LangGraphClient struct {
	Assistants *client.AssistantsClient
	Threads    *client.ThreadsClient
	Runs       *client.RunsClient
	Crons      *client.CronsClient
	Store      *client.StoreClient
}

func newLangGraphClient(httpClient *http.HttpClient) *LangGraphClient {
	return &LangGraphClient{
		Assistants: client.NewAssistantsClient(httpClient),
		Threads:    client.NewThreadsClient(httpClient),
		Runs:       client.NewRunsClient(httpClient),
		Crons:      client.NewCronsClient(httpClient),
		Store:      client.NewStoreClient(httpClient),
	}
}

func getApiKey(apiKey string) string {
	if apiKey != "" {
		return apiKey
	}

	prefixes := []string{"LANGGRAPH", "LANGSMITH", "LANGCHAIN"}
	for _, prefix := range prefixes {
		if env := os.Getenv(fmt.Sprintf("%s_API_KEY", prefix)); env != "" {
			return strings.TrimSpace(env)
		}
	}
	return ""
}

func getHeaders(apiKey string, customHeaders map[string]string) map[string]string {
	for _, header := range RESERVED_HEADERS {
		if _, exists := customHeaders[header]; exists {
			panic(fmt.Sprintf("Cannot set reserved header '%s'", header))
		}
	}

	headers := map[string]string{
		"User-Agent": fmt.Sprintf("langgraph-sdk-go/%s", Version),
	}
	maps.Copy(headers, customHeaders)

	apiKey = getApiKey(apiKey)
	if apiKey != "" {
		headers["x-api-key"] = apiKey
	}

	return headers
}

func GetClient(url string, apiKey string, headers map[string]string) *LangGraphClient {
	if url == "" {
		url = "http://localhost:2024"
	}

	transport := &http_client.Transport{
		Proxy:               http_client.ProxyFromEnvironment,
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	httpWrapper := http.NewHttpClient(
		url,
		getHeaders(apiKey, headers),
		300*time.Second,
		transport,
	)

	maxRetries := 5
	retryInterval := 3 * time.Second
	var lastErr error

	for i := range maxRetries {
		err := httpWrapper.CheckConnection()
		if err == nil {
			break
		}

		lastErr = err
		if i < maxRetries-1 {
			time.Sleep(retryInterval)
		}
	}

	if lastErr != nil {
		panic(fmt.Sprintf("Failed to connect after %d attempts: %v", maxRetries, lastErr))
	}

	return newLangGraphClient(httpWrapper)
}
