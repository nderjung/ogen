package sse_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ogen-go/ogen/sse"
)

func TestEventStreamIntegration(t *testing.T) {
	// Create a simple event payload
	type Event struct {
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
	}

	// Set up a test server that sends SSE events
	eventsReceived := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		// Send 3 events
		for i := 0; i < 3; i++ {
			event := Event{
				Message:   fmt.Sprintf("Event %d", i),
				Timestamp: time.Now().Unix(),
			}

			data, _ := json.Marshal(event)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
			time.Sleep(50 * time.Millisecond)
		}
	}))
	defer server.Close()

	// Make a client request
	req, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify content type
	require.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))

	// Use the SSE Reader to process events
	reader := sse.NewReader(resp.Body)

	// Process events similar to how generated code would
	for i := 0; i < 3; i++ {
		event, err := reader.ReadEvent()
		require.NoError(t, err)

		var decoded Event
		err = json.Unmarshal(event.Data, &decoded)
		require.NoError(t, err)

		require.Contains(t, decoded.Message, fmt.Sprintf("Event %d", i))
		require.NotZero(t, decoded.Timestamp)

		eventsReceived++
	}

	// After all events are processed, next read should block or return EOF
	// since we have a 1 second timeout, this should exit with EOF or context deadline
	_, err = reader.ReadEvent()
	require.True(t, err == io.EOF || errors.Is(err, context.DeadlineExceeded))

	// Verify we received all events
	require.Equal(t, 3, eventsReceived)
}

func ExampleReader() {
	// This example demonstrates how to use the SSE Reader directly

	// Create a reader from a string containing SSE data
	input := "data: {\"message\":\"hello\",\"count\":1}\n\n" +
		"data: {\"message\":\"world\",\"count\":2}\n\n"
	reader := sse.NewReader(strings.NewReader(input))

	// Read all events
	for {
		event, err := reader.ReadEvent()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		fmt.Printf("Received event data: %s\n", event.Data)
	}

	// Output:
	// Received event data: {"message":"hello","count":1}
	// Received event data: {"message":"world","count":2}
}
