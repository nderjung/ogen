package sse_test

import (
	"io"
	"strings"
	"testing"

	"github.com/ogen-go/ogen/sse"
	"github.com/stretchr/testify/require"
)

func TestSSEReader(t *testing.T) {
	t.Run("SingleEvent", func(t *testing.T) {
		input := "data: {\"message\":\"hello world\",\"count\":42}\n\n"
		reader := sse.NewReader(strings.NewReader(input))

		event, err := reader.ReadEvent()
		require.NoError(t, err)
		require.NotNil(t, event)

		require.Equal(t, []byte("{\"message\":\"hello world\",\"count\":42}"), event.Data)
		require.Nil(t, event.ID)
		require.Nil(t, event.Event)
		require.Nil(t, event.Retry)
	})

	t.Run("MultipleEvents", func(t *testing.T) {
		input := "data: event1\n\ndata: event2\n\ndata: event3\n\n"
		reader := sse.NewReader(strings.NewReader(input))

		// Read first event
		event1, err := reader.ReadEvent()
		require.NoError(t, err)
		require.Equal(t, []byte("event1"), event1.Data)

		// Read second event
		event2, err := reader.ReadEvent()
		require.NoError(t, err)
		require.Equal(t, []byte("event2"), event2.Data)

		// Read third event
		event3, err := reader.ReadEvent()
		require.NoError(t, err)
		require.Equal(t, []byte("event3"), event3.Data)

		// Should return EOF now
		_, err = reader.ReadEvent()
		require.Equal(t, io.EOF, err)
	})

	t.Run("EmptyLines", func(t *testing.T) {
		input := "\n\ndata: event\n\n\n\n"
		reader := sse.NewReader(strings.NewReader(input))

		event, err := reader.ReadEvent()
		require.NoError(t, err)
		require.Equal(t, []byte("event"), event.Data)

		_, err = reader.ReadEvent()
		require.Equal(t, io.EOF, err)
	})

	t.Run("CommentLines", func(t *testing.T) {
		input := ": this is a comment\ndata: actual event\n\n"
		reader := sse.NewReader(strings.NewReader(input))

		event, err := reader.ReadEvent()
		require.NoError(t, err)
		require.Equal(t, []byte("actual event"), event.Data)
	})

	t.Run("AllFieldTypes", func(t *testing.T) {
		input := "event: update\nid: 123\ndata: event data\nretry: 3000\n\n"
		reader := sse.NewReader(strings.NewReader(input))

		event, err := reader.ReadEvent()
		require.NoError(t, err)
		require.Equal(t, []byte("event data"), event.Data)
		require.Equal(t, []byte("update"), event.Event)
		require.Equal(t, []byte("123"), event.ID)
		require.Equal(t, []byte("3000"), event.Retry)
	})

	t.Run("MultilineData", func(t *testing.T) {
		input := "data: line1\ndata: line2\ndata: line3\n\n"
		reader := sse.NewReader(strings.NewReader(input))

		event, err := reader.ReadEvent()
		require.NoError(t, err)
		require.Equal(t, []byte("line1\nline2\nline3"), event.Data)
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		// Properly formatted data after invalid line
		input := "data without colon\ndata: valid\n\n"
		reader := sse.NewReader(strings.NewReader(input))

		event, err := reader.ReadEvent()
		require.NoError(t, err)
		require.Equal(t, []byte("valid"), event.Data)
	})

	t.Run("DifferentLineEndings", func(t *testing.T) {
		// Test with \r\n line endings
		input := "data: windows\r\n\r\n"
		reader := sse.NewReader(strings.NewReader(input))

		event, err := reader.ReadEvent()
		require.NoError(t, err)
		require.Equal(t, []byte("windows"), event.Data)
	})

	t.Run("KeepAlive", func(t *testing.T) {
		// Test with keep-alive comments followed by data
		input := ": keep-alive\ndata: real event\n\n"
		reader := sse.NewReader(strings.NewReader(input))

		event, err := reader.ReadEvent()
		require.NoError(t, err)
		require.Equal(t, []byte("real event"), event.Data)
	})
}
