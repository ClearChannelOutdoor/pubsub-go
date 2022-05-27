package pubsub_go

type Message struct {
	Attributes map[string]string
	Message    interface{}
	Topic      string
}
