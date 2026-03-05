package intake

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"time"
)

var (
	errInvalidJSON         = errors.New("invalid JSON")
	errMultipleJSONObjects = errors.New("request body must contain only one JSON object")
	errMissingTimestamp    = errors.New("timestamp is required")
	errTimestampNotString  = errors.New("timestamp must be a string")
	errTimestampNotRFC3339 = errors.New("timestamp must be RFC3339")
)

func isJSONContentType(contentType string) bool {
	if contentType == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return mediaType == "application/json"
}

func decodeJSONBody(r io.Reader) (map[string]any, error) {
	var body map[string]any
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&body); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, io.EOF
		}
		return nil, errInvalidJSON
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, errMultipleJSONObjects
	}

	tsRaw, ok := body["timestamp"]
	if !ok {
		return nil, errMissingTimestamp
	}
	ts, ok := tsRaw.(string)
	if !ok {
		return nil, errTimestampNotString
	}
	if _, err := time.Parse(time.RFC3339, ts); err != nil {
		return nil, errTimestampNotRFC3339
	}

	return body, nil
}
