// Package pubsub_go TODO: Description
package pubsub_go

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"google.golang.org/api/option"
)

// PubSub provides a wrapper client for Google Cloud's PubSub for
// publishing messages to a topic and receiving messages from a subscription.
type PubSub struct {
	client   *pubsub.Client
	settings settings
}

type settings struct {
	publish PublishSettings
	receive ReceiveSettings
}

// Config provides the information needed to securely connect to Google Cloud's PubSub
// and to configure any publishing and subscription options.
type Config struct {
	ProjectID              string
	IsLocal                bool
	ServiceAccountFilePath string
	PublishSettings        PublishSettings
	ReceiveSettings        ReceiveSettings
}

// PublishSettings is an extension of Google PubSub's PublishSettings that
// enables further configuration for publishing messages to a topic.
type PublishSettings struct {
	Settings pubsub.PublishSettings
	// TODO: Add more config like auto creating a topic if it doesn't exist
}

// ReceiveSettings is an extension of Google PubSub's ReceiveSettings that
// enables further configuration for receiving messages from a subscription.
type ReceiveSettings struct {
	Settings pubsub.ReceiveSettings
}

// NewPubSub creates a new PubSub client with the provided Config.
func NewPubSub(c Config) (*PubSub, error) {
	client, err := newClient(c.ProjectID, c.IsLocal, c.ServiceAccountFilePath)
	if err != nil {
		return nil, err
	}

	return &PubSub{
		client: client,
		settings: settings{
			publish: c.PublishSettings,
			receive: c.ReceiveSettings,
		},
	}, nil
}

// Publish sends a message to a topic along with any attributes that were provided.
func (ps *PubSub) Publish(m Message) error {
	topic := ps.client.Topic(m.Topic)
	topic.PublishSettings = ps.settings.publish.Settings

	data, err := json.Marshal(m.Message)
	if err != nil {
		return err
	}

	ctx := context.Background()
	result := topic.Publish(ctx, &pubsub.Message{
		Data:       data,
		Attributes: m.Attributes,
	})

	_, err = result.Get(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Receive subscribes to a topic via the subscription id and passes messages back to the caller
// through the channel.
func (ps *PubSub) Receive(subscription string, messages chan<- *pubsub.Message) error {
	ctx := context.Background()
	s := ps.client.Subscription(subscription)

	return s.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		messages <- m
	})
}

func newClient(projectID string, isLocal bool, settingsPath string) (*pubsub.Client, error) {
	ctx := context.Background()

	if isLocal {
		return pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(settingsPath))
	}

	return pubsub.NewClient(ctx, projectID)
}
