package pubsub_go

// Message holds the data needed for publishing a message to PubSub.
type Message struct {
	Attributes map[string]string
	Message    interface{}
	Topic      string
}
