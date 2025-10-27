package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"to-tcp/internal/headers"
	"to-tcp/internal/request"
	"to-tcp/internal/response"
	"to-tcp/internal/server"
)

const PORT int = 42069

const CHUNK_SIZE int = 1024

func badReqResp() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}
func internalErrResp() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}
func okResp() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func main() {
	s, err := server.Serve(
		PORT,
		func(w *response.Writer, req *request.Request) {
			h := response.GetDefaultHeaders(0)
			body := okResp()
			status := response.StatusOK

			switch {
			case req.RequestLine.RequestTarget == "/badRequest":
				body = badReqResp()
				status = response.StatusBadRequest
			case req.RequestLine.RequestTarget == "/internalErr":
				body = internalErrResp()
				status = response.StatusInternalServerErr
			case strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream"):
				target := req.RequestLine.RequestTarget

				httpbinRes, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
				if err != nil {
					body = internalErrResp()
					status = response.StatusInternalServerErr
				} else {
					// chunked response
					w.WriteStatusLine(status)

					h.Delete("Content-Length")
					h.Set("Transfer-Encoding", "chunked")
					h.Replace("Content-Type", "text/plain")
					h.Set("Trailer", "X-Content-SHA256")
					h.Set("Trailer", "X-Content-Length")

					w.WriteHeaders(h)

					fullBody := []byte{}
					chunk := make([]byte, CHUNK_SIZE)
					for {
						n, err := httpbinRes.Body.Read(chunk)
						if n > 0 {
							w.WriteChunkedBody(chunk, n)
						}
						if err != nil {
							break
						}

						fullBody = append(fullBody, chunk[:n]...)
					}
					w.WriteChunkedBodyDone()

					sha256Checksum := sha256.Sum256(fullBody)
					hexChecksum := hex.EncodeToString(sha256Checksum[:])

					trailers := headers.NewHeaders()
					trailers.Set("X-Content-SHA256", hexChecksum)
					trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
					w.WriteHeaders(trailers)

					return
				}
			}

			h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
			h.Replace("Content-Type", "text/html")
			w.WriteStatusLine(status)
			w.WriteHeaders(h)
			w.WriteBody(body)
		},
	)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
