package ir

import "strings"

// ContentType is a Content-Type header value.
type ContentType string

func (t ContentType) Mask() bool { return strings.ContainsRune(string(t), '*') }

func (t ContentType) String() string { return string(t) }

// Encoding of body.
type Encoding string

const (
	// EncodingJSON is Encoding for json.
	EncodingJSON Encoding = "application/json"
	// EncodingFormURLEncoded is Encoding for URL-encoded form.
	EncodingFormURLEncoded Encoding = "application/x-www-form-urlencoded"
	// EncodingMultipart is Encoding for multipart form.
	EncodingMultipart Encoding = "multipart/form-data"
	// EncodingOctetStream is Encoding for binary.
	EncodingOctetStream Encoding = "application/octet-stream"
	// EncodingTextPlain is Encoding for text.
	EncodingTextPlain Encoding = "text/plain"
	// EncodingEventStream is Encoding for server-sent events.
	EncodingEventStream Encoding = "text/event-stream"
)

func (t Encoding) String() string { return string(t) }

func (t Encoding) JSON() bool { return t == EncodingJSON }

func (t Encoding) FormURLEncoded() bool { return t == EncodingFormURLEncoded }

func (t Encoding) MultipartForm() bool { return t == EncodingMultipart }

func (t Encoding) OctetStream() bool { return t == EncodingOctetStream }

func (t Encoding) TextPlain() bool { return t == EncodingTextPlain }

func (t Encoding) EventStream() bool { return t == EncodingEventStream }

type Media struct {
	// Encoding is the parsed content type used for encoding, but not for header value.
	Encoding Encoding
	// Type is response or request type.
	Type *Type

	// JSONStreaming indicates that the JSON media should be streamed.
	JSONStreaming bool

	// EventStreaming indicates that the media should be handled as Server-Sent Events.
	EventStreaming bool
}
