package http

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpClient_Get(t *testing.T) {
	client := NewHttpClient("http://localhost:2024", nil, 0, nil)

	_, err := client.Get(context.Background(), "/test-path", nil, nil, nil)

	assert.NoError(t, err, "Expected no error when sending GET request")
}
