# Trade Republic WebSocket Client

A Go WebSocket client for connecting to Trade Republic's real-time data streams. This package provides a type-safe, concurrent WebSocket client with automatic message handling and subscription management.

## Installation

```bash
go get github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic
```

## Features

- **Type-safe WebSocket communication** with generated structs
- **Concurrent subscription handling** with goroutines
- **Automatic message parsing** and state management
- **Publisher-subscriber pattern** for message distribution
- **Context-based cancellation** support
- **Thread-safe operations** with mutex protection
- **Mock support** for testing

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Create publisher for message distribution
    publisher := traderepublic.NewPublisher()
    
    // Create WebSocket client
    client := traderepublic.NewClient(publisher, ctx)
    
    // Connect to WebSocket
    if err := client.Connect(); err != nil {
        log.Fatal("Failed to connect:", err)
    }
    defer client.Close()
    
    // Create subscription request
    subRequest := traderepublic.WsSubRequestJson{
        Type: "timeline",
        // Add other required fields
    }
    
    // Subscribe to data stream
    dataChan, err := client.Subscribe(subRequest)
    if err != nil {
        log.Fatal("Failed to subscribe:", err)
    }
    
    // Read messages
    for {
        select {
        case data := <-dataChan:
            fmt.Printf("Received: %s\n", string(data))
        case <-ctx.Done():
            return
        }
    }
}
```

## Interfaces

### WSClientInterface

The main client interface:

```go
type WSClientInterface interface {
    // Connect establishes WebSocket connection
    Connect() error
    
    // Close terminates the WebSocket connection  
    Close() error
    
    // Subscribe creates a subscription and returns data channel
    Subscribe(data WsSubRequestJson) (<-chan []byte, error)
}
```

### PublisherInterface

Message distribution interface:

```go
type PublisherInterface interface {
    // Subscribe creates a subscription channel
    Subscribe(subID string) <-chan []byte
    
    // Publish sends data to subscribers
    Publish(data []byte, subID string)
    
    // Close closes a subscription
    Close(subID string)
}
```

## Data Types

### WebSocket Request

```go
type WsSubRequestJson struct {
    Type    string `json:"type"`
    ID      string `json:"id,omitempty"`
    // Additional fields based on subscription type
}
```

### WebSocket Response

```go
type WsResponseJson struct {
    State WsResponseJsonState `json:"state"`
    Body  interface{}         `json:"body,omitempty"`
    // Response data
}
```

### Message States

- `StateData` (`"A"`): Complete data message
- `StateContinue` (`"C"`): Partial data, more coming
- `StateError` (`"E"`): Error occurred

## Error Handling

The client defines specific error types:

```go
var (
    ErrNotConnected     = errors.New("websocket not connected")
    ErrConnectionClosed = errors.New("websocket connection closed")
    ErrAuthRequired     = errors.New("authentication required")
)
```

Example error handling:

```go
dataChan, err := client.Subscribe(subRequest)
if err != nil {
    switch {
    case errors.Is(err, traderepublic.ErrNotConnected):
        // Handle connection error
        log.Println("Not connected to WebSocket")
    case errors.Is(err, traderepublic.ErrAuthRequired):
        // Handle authentication error
        log.Println("Authentication required")
    default:
        // Handle other errors
        log.Printf("Subscription failed: %v", err)
    }
}
```

## Advanced Usage

### Custom Publisher

```go
type MyPublisher struct {
    channels map[string]chan []byte
    mu       sync.RWMutex
}

func NewMyPublisher() *MyPublisher {
    return &MyPublisher{
        channels: make(map[string]chan []byte),
    }
}

func (p *MyPublisher) Subscribe(subID string) <-chan []byte {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    ch := make(chan []byte, 1000) // Large buffer
    p.channels[subID] = ch
    return ch
}

func (p *MyPublisher) Publish(data []byte, subID string) {
    p.mu.RLock()
    ch, exists := p.channels[subID]
    p.mu.RUnlock()
    
    if exists {
        select {
        case ch <- data:
        default:
            // Handle full channel
            log.Printf("Channel full for subscription %s", subID)
        }
    }
}

func (p *MyPublisher) Close(subID string) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if ch, exists := p.channels[subID]; exists {
        close(ch)
        delete(p.channels, subID)
    }
}
```

### Multiple Subscriptions

```go
// Subscribe to multiple data streams
timelineReq := traderepublic.WsSubRequestJson{Type: "timeline"}
instrumentReq := traderepublic.WsSubRequestJson{Type: "instrument"}

timelineChan, err := client.Subscribe(timelineReq)
if err != nil {
    log.Fatal(err)
}

instrumentChan, err := client.Subscribe(instrumentReq)
if err != nil {
    log.Fatal(err)
}

// Handle multiple channels
go func() {
    for data := range timelineChan {
        // Process timeline data
        handleTimelineData(data)
    }
}()

go func() {
    for data := range instrumentChan {
        // Process instrument data
        handleInstrumentData(data)
    }
}()
```

### Connection Management

```go
func createResilientClient(ctx context.Context) traderepublic.WSClientInterface {
    publisher := traderepublic.NewPublisher()
    
    for {
        client := traderepublic.NewClient(publisher, ctx)
        
        if err := client.Connect(); err != nil {
            log.Printf("Connection failed: %v, retrying...", err)
            time.Sleep(5 * time.Second)
            continue
        }
        
        return client
    }
}
```

## Testing

The package includes generated mocks for testing:

```go
func TestMyService(t *testing.T) {
    // Create mock client
    mockClient := traderepublic.NewMockWSClientInterface(gomock.NewController(t))
    
    // Set expectations
    dataChan := make(chan []byte, 1)
    mockClient.EXPECT().
        Subscribe(gomock.Any()).
        Return(dataChan, nil)
    
    // Test your code
    service := NewMyService(mockClient)
    err := service.StartSubscription()
    assert.NoError(t, err)
}
```

## Configuration

### WebSocket Endpoint

The client connects to:
```
wss://api.traderepublic.com/v1/websocket
```

### Message Format

Messages follow this pattern:
```
<messageId> <messageType> <data>
```

Where:
- `messageId`: Unique identifier
- `messageType`: `sub`, `unsub`, or response state
- `data`: JSON payload

## Logging

The client uses structured logging with `slog`:

```go
import "log/slog"

// Configure logging level
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))
slog.SetDefault(logger)
```

## Thread Safety

All client operations are thread-safe:
- Multiple goroutines can call `Subscribe()` concurrently
- Message reading and distribution happen in separate goroutines
- Internal state is protected by mutexes

## Context Support

The client respects context cancellation:

```go
// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

client := traderepublic.NewClient(publisher, ctx)

// Client will automatically close when context is cancelled
```

## Generated Code

This package uses code generation for:
- **Type definitions** from JSON schemas
- **Mock interfaces** for testing
- **OpenAPI client** code

To regenerate code:
```bash
go generate ./...
```

## Contributing

1. Ensure all tests pass: `go test ./...`
2. Update schemas in `schemas/` directory if needed
3. Run code generation: `go generate ./...`
4. Follow Go best practices and add tests for new features

## License

See the main repository for license information.

## Examples

Check the repository for complete examples and integration tests.