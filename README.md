# pubsub-go
pubsub-go is a Go module that makes it easy to publish and receive messages from GCP's Pub/Sub service.

## Usage
### Create a new PubSub
```go
import pubsub_go "github.com/clearchanneloutdoor/pubsub-go"

func main() {
    ctx := context.Background()
    config := pubsub_go.Config{
        ProjectID:              "projectID",
        IsLocal:                true,
        ServiceAccountFilePath: "./path/to/settings.json",
    }
    
    ps, err := pubsub_go.NewClient(ctx, config)	
    if err != nil {
        // Handle error
    }
}
```

### Publish Message
```go
func main() {
    // Initialize new pubsub-go PubSub
	
    message := pubsub_go.Message{
        Attributes: nil,
        Message:    "Hello World",
        Topic:      "topic",
    }
    
    if err := ps.Publish(message); err != nil {
        // Handle error
    }
}
```

### Receive Message
```go
func main() {
    // Initialize new pubsub-go PubSub
	
    messages := make(chan *pubsub.Message)
	
    go func() {
        for {
            msg := <-messages
			
            // Process message
            fmt.Printf("Got message: %s\n", string(msg.Data))
			
            // Acknowledge message
            msg.Ack()
        }
    }()
    
    if err := ps.Receive("subscription", messages); err != nil {
        // Handle error
    }
}
```