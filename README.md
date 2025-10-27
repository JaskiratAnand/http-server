# HTTP 1.1 Server - Golang

This project is a from-scratch implementation of an HTTP/1.1 server built using Go's standard library. It's designed to demonstrate the fundamentals of parsing raw TCP streams, constructing HTTP requests from them, and sending back valid HTTP responses.

## Working Features

*   **Raw HTTP/1.1 Parsing**: The server can parse the request line, headers, and body from an incoming TCP stream.
*   **Static Responses**: Serves basic HTML content for different endpoints.
*   **Endpoint Routing**: A simple `switch` statement in the handler routes requests to different logic based on the request target.
*   **Chunked Transfer Encoding**: The server can act as a proxy to stream a response from another service (`httpbin.org` in this case) using chunked encoding. This is handled efficiently using Go's `io.Copy` and a custom `io.Writer` implementation.
*   **Transfering Binary Data**: The server listens for `SIGINT` and `SIGTERM` signals to shut down gracefully.
*   **Graceful Shutdown**: The server listens for `SIGINT` and `SIGTERM` signals to shut down gracefully.

## Learnings

Building this server provided several key insights:

*   **HTTP/1.1 Protocol**: I gained a deeper understanding of the HTTP/1.1 specification, particularly the structure of requests and responses, the importance of `CRLF` (`\r\n`), and the mechanism of chunked transfer encoding.
*   **TCP Sockets**: This project required working directly with TCP sockets (`net.Listener` and `net.Conn`), providing a practical look at how application-level protocols are built on top of TCP.
*   **Go's `io` Package**: The `io.Reader` and `io.Writer` interfaces are incredibly powerful. By creating a custom `ChunkedResponseWriter` that implements `io.Writer`, I was able to leverage `io.Copy` to create an elegant and efficient streaming proxy with minimal code.
*   **API Design**: Designing the `response.Writer` and `request.Request` types helped in understanding how to create abstractions that make the higher-level application logic (the server handler) cleaner and easier to manage.

## How to Create Your Own Server Handler

The server is designed to be extensible. You can easily create your own handlers to add new endpoints and logic.

The core of the server is the `server.Serve` function, which accepts a port and a handler function. The handler has the following signature:

```go
func(w *response.Writer, req *request.Request)
```

-   `w *response.Writer`: This is your tool for writing the HTTP response back to the client. You can use it to set the status line, write headers, and send the body.
-   `req *request.Request`: This struct contains the parsed information from the client's request, including the request line, headers, and body.

### Example Handler

Here is a simple example of how you could add a new endpoint that returns a JSON response:

```go
// Inside the main handler function in cmd/httpserver/main.go

func(w *response.Writer, req *request.Request) {
    h := response.GetDefaultHeaders(0)
    body := okResp()
    status := response.StatusOK

    switch {
    // ... other cases

    case req.RequestLine.RequestTarget == "/api/health":
        // This is our new handler logic
        status = response.StatusOK
        body = []byte(`{"status": "ok"}`)
        h.Replace("Content-Type", "application/json")
        h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))

        w.WriteStatusLine(status)
        w.WriteHeaders(h)
        w.WriteBody(body)
        return // Return early to avoid the default response logic

    // ... other cases
    }

    // Default response logic
    h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
    h.Replace("Content-Type", "text/html")
    w.WriteStatusLine(status)
    w.WriteHeaders(h)
    w.WriteBody(body)
}
```

## How to Run

You can run either of the applications using `go run`. Make sure you are in the root directory of the project.

**Note:** You can only run one of the applications at a time, as they both bind to the same port (`42069`).

### Running the HTTP Server

1.  Start the server:
    ```sh
    go run ./cmd/httpserver
    ```
    You should see the output: `Server started on port 42069`.

2.  In a separate terminal, send requests to the server using `curl`:
    ```sh
    # Get a 200 OK response
    curl http://localhost:42069/

    # Get a 400 Bad Request response
    curl http://localhost:42069/badRequest

    # Get a 500 Internal Server Error response
    curl http://localhost:42069/internalErr

    # Stream a chunked response from httpbin
    curl http://localhost:42069/httpbin/stream/20

    # Stream a video
    curl http://localhost:42069/video
    ```

### Running the TCP Listener (`cmd/tcplistner`)

This application is a simple TCP listener that demonstrates the request parsing logic. It accepts a connection, parses the HTTP request, prints it to the console, and then closes the connection.

1.  Start the listener:
    ```sh
    go run ./cmd/tcplistner
    ```

2.  In a separate terminal, send a request using `curl`.
    ```sh
    curl -X POST -d "Hello, this is the body" http://localhost:42069/some/path
    ```
    The listener's output will show the fully parsed request.

## Testing

The project includes unit tests for the low-level packages to ensure they are working correctly. You can run all the tests using the following command from the root of the project:

```sh
go test ./...
```

This will run tests for:

*   **`internal/headers`**: Tests the parsing and manipulation of HTTP headers.
*   **`internal/request`**: Tests the parsing of the HTTP request line.
