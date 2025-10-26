package response

import (
	"fmt"
	"io"
	"to-tcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                StatusCode = 200
	StatusCreated           StatusCode = 201
	StatusBadRequest        StatusCode = 400
	StatusUnauthorized      StatusCode = 401
	StatusForbidden         StatusCode = 403
	StatusNotFound          StatusCode = 404
	StatusInternalServerErr StatusCode = 500
)

// type WriteStatus uint8
// const (
// 	WriteStatusLine WriteStatus = 1
// 	WriteHeaders    WriteStatus = 2
// 	WriteBody       WriteStatus = 3
// )

type Writer struct {
	writer io.Writer
	// writeStatus WriteStatus
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOK:
		statusLine = []byte("HTTP/1.1 200 OK")
	case StatusCreated:
		statusLine = []byte("HTTP/1.1 201 Created")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request")
	case StatusUnauthorized:
		statusLine = []byte("HTTP/1.1 401 Unauthorized")
	case StatusForbidden:
		statusLine = []byte("HTTP/1.1 403 Forbidden")
	case StatusNotFound:
		statusLine = []byte("HTTP/1.1 404 Not Found")
	case StatusInternalServerErr:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error")
	default:
		return fmt.Errorf("unrecognized error code")
	}

	statusLine = fmt.Append(statusLine, "\r\n")
	_, err := w.writer.Write(statusLine)
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Content-Type", "text/plain")
	h.Set("Connection", "close")

	return h
}

func (w *Writer) WriteHeaders(h *headers.Headers) error {
	b := []byte{}
	h.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Append(b, "\r\n")

	_, err := w.writer.Write(b)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)

	return n, err
}
