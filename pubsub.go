// Package pubsub_go TODO: Description
package pubsub_go

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
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
	IsLocal                bool
	ProjectID              string
	PublishSettings        PublishSettings
	ReceiveSettings        ReceiveSettings
	ServiceAccountFilePath string
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

// SubscriptionConfig is an extension of Google PubSub's SubscriptionConfig that
// enables further configuration for a subscription of a topic
type SubscriptionConfig struct {
	Settings pubsub.SubscriptionConfig
}

// TopicConfig is an extension of Google PubSub's TopicConfig that
// enables further configuration for a topic
type TopicConfig struct {
	Settings pubsub.TopicConfig
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

// Create Subscriptions for a Topic based on a map of Subscription Name and Filter
func (ps *PubSub) CreateSubscriptions(tid string, sids map[string]string, cfg *SubscriptionConfig) error {
	ctx := context.Background()

	// find the topic by tid first
	topic := ps.client.Topic(tid)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return err
	}
	// if the topic does not exist, return an error
	if !exists {
		return fmt.Errorf("topic %s does not exist", tid)
	}

	cfg.Settings.Topic = topic
	// let's create subscriptions with optional filters for the topic
	for sid, flt := range sids {
		// check if the subscription exists or not
		sub := ps.client.Subscription(sid)
		exists, err := sub.Exists(ctx)
		if err != nil {
			return err
		}
		// if the sub already exists, skip it
		if exists {
			continue
		}
		// only if the filter is provided in the map[string]string to include the Filter config
		if flt != "" {
			cfg.Settings.Filter = flt
		}
		// creating a subscription
		_, err = ps.client.CreateSubscription(ctx, sid, cfg.Settings)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create a Topic in Google PubSub if not exist
func (ps *PubSub) CreateTopic(tid string, cfg *TopicConfig) error {
	ctx := context.Background()

	topic := ps.client.Topic(tid)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = ps.client.CreateTopicWithConfig(ctx, tid, &cfg.Settings)
	if err != nil {
		return err
	}
	return nil
}

// Publish sends a message to a topic along with any attributes that were provided.
func (ps *PubSub) Publish(m Message) error {
	topic := ps.client.Topic(m.Topic)
	topic.PublishSettings = ps.settings.publish.Settings

	data, err := json.Marshal(m.Message)
	if err != nil {
		return err
	}

	if m.Attributes != nil {
		m.Attributes["OriginatedAt"] = fmt.Sprintf("%v", time.Now().Unix())
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
	s.ReceiveSettings = pubsub.DefaultReceiveSettings

	return s.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		messages <- m
	})
}

func (ps *PubSub) ReceiveWithSettings(subscription string, rs ReceiveSettings, messages chan<- *pubsub.Message) error {
	ctx := context.Background()
	s := ps.client.Subscription(subscription)

	s.ReceiveSettings = pubsub.DefaultReceiveSettings

	// override the default settings with the provided settings
	if rs.Settings.MaxExtension != 0 {
		s.ReceiveSettings.MaxExtension = rs.Settings.MaxExtension
	}
	if rs.Settings.MaxExtensionPeriod != 0 {
		s.ReceiveSettings.MaxExtensionPeriod = rs.Settings.MaxExtensionPeriod
	}
	if rs.Settings.MinExtensionPeriod != 0 {
		s.ReceiveSettings.MinExtensionPeriod = rs.Settings.MinExtensionPeriod
	}
	if rs.Settings.MaxOutstandingMessages != 0 {
		s.ReceiveSettings.MaxOutstandingMessages = rs.Settings.MaxOutstandingMessages
	}
	if rs.Settings.MaxOutstandingBytes != 0 {
		s.ReceiveSettings.MaxOutstandingBytes = rs.Settings.MaxOutstandingBytes
	}
	if rs.Settings.NumGoroutines != 0 {
		s.ReceiveSettings.NumGoroutines = rs.Settings.NumGoroutines
	}

	// always false to allow the library to use streamingpull-api
	s.ReceiveSettings.Synchronous = false

	return s.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		messages <- m
	})
}

func newClient(projectID string, isLocal bool, settingsPath string) (*pubsub.Client, error) {
	ctx := context.Background()

	// if isLocal is true or is settingsPath is empty, then use credentials
	useCreds := isLocal || settingsPath != ""
	if useCreds {
		return pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(settingsPath))
	}

	return pubsub.NewClient(ctx, projectID)
}
