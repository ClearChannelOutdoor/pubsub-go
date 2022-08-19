# pubsub-go
pubsub-go is a Go module that makes it easy to publish and receive messages from GCP's Pub/Sub service.

## Usage
### Create a new PubSub
```go
import pubsub_go "github.com/clearchanneloutdoor/pubsub-go"

func main() {
    ctx := context.Background()
    config := pubsub_go.Config{
        ProjectID:                      "projectID",
        IsLocal:                        true,
        ServiceAccountFilePath:         "./path/to/settings.json",
        SubscriptionEnableOrdering:     true,
        SubscriptionRetainAckedMessages: false,
        TopicRetentionDays:             3,
    }
    
    ps, err := pubsub_go.NewClient(ctx, config)	
    if err != nil {
        // Handle error
    }
}
```

### Create a Topic
```go
func main() {
    // Initialize new pubsub-go PubSub

    if err := ps.CreateTopic("topic"); err != nil {
        // Handle error
    }
    
    // send a message to the topic}
```

### Create Multiple Subscriptions with Filters for a Topic
```go
func main() {
    // Initialize new pubsub-go PubSub

    // Define a map for a list of Subscription Names and their (optional) filter definitions
    subs := map[string]string{
        "topic-sub": "",
        "topic-sub-ca": "attributes.region = \"CA\"",
        "topic-sub-nv": "attributes.region = \"NV\"",
        "topic-sub-tx": "attributes.region = \"TX\"",
    }
    
    tid := "topic"
    // create subscriptions
    if err := ps.CreateSubscriptions(tid, subs); err != nil {
        // handle error
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