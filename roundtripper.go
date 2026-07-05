package log

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LoggingRoundTripper is an HTTP RoundTripper that logs requests and responses
type LoggingRoundTripper struct {
	Transport http.RoundTripper
	Logger    Logger
	LogBodies bool
}

// RoundTrip executes an HTTP request and logs the result
func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	logEntry := lrt.Logger.WithContext(req.Context())

	var reqBody []byte
	if lrt.LogBodies && req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	resp, err := lrt.Transport.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		msg := fmt.Sprintf("HTTP request failed: %s %s, error: %s, duration: %dms",
			req.Method, req.URL.String(), err.Error(), duration.Milliseconds())
		if lrt.LogBodies && len(reqBody) > 0 {
			msg += fmt.Sprintf(", request: %s", string(reqBody))
		}
		logEntry.Error(msg)
		return nil, err
	}

	var respBody []byte
	if lrt.LogBodies && resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	}

	if lrt.LogBodies {
		logEntry.Infof("HTTP request completed: %s %s, duration: %dms, request: %s, response: %s",
			req.Method, req.URL.String(), duration.Milliseconds(), string(reqBody), string(respBody))
	} else {
		logEntry.Infof("HTTP request completed: %s %s, duration: %dms",
			req.Method, req.URL.String(), duration.Milliseconds())
	}
	return resp, nil
}
