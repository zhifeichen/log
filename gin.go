package log

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func newResponseBodyWriter(w gin.ResponseWriter) *responseBodyWriter {
	return &responseBodyWriter{
		ResponseWriter: w,
		body:           bytes.NewBuffer(nil),
	}
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// GinLogger returns a Gin middleware that logs HTTP requests with traceID
func GinLogger(l Logger, logBodies bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := l.SetContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()

		var reqBody []byte
		if logBodies && c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body.Close()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		rw := newResponseBodyWriter(c.Writer)
		c.Writer = rw

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		url := c.Request.URL.String()

		logEntry := l.WithContext(ctx)
		if status >= 400 {
			msg := fmt.Sprintf("request failed: %d %s %s, duration: %dms", status, method, url, duration.Milliseconds())
			if logBodies {
				msg += fmt.Sprintf(", request: %s, response: %s", string(reqBody), rw.body.String())
			}
			logEntry.Error(msg)
		} else {
			msg := fmt.Sprintf("request completed: %d %s %s, duration: %dms", status, method, url, duration.Milliseconds())
			if logBodies {
				msg += fmt.Sprintf(", request: %s, response: %s", string(reqBody), rw.body.String())
			}
			logEntry.Info(msg)
		}
	}
}
