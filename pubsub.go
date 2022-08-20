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
	publish      PublishSettings
	receive      ReceiveSettings
	subscription SubscriptionSettings
	topic        TopicSettings
}

// Config provides the information needed to securely connect to Google Cloud's PubSub
// and to configure any publishing and subscription options.
type Config struct {
	IsLocal                         bool
	ProjectID                       string
	PublishSettings                 PublishSettings
	ReceiveSettings                 ReceiveSettings
	ServiceAccountFilePath          string
	SubscriptionEnableOrdering      bool
	SubscriptionRetainAckedMessages bool
	TopicRetentionDays              int
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

// SubscriptionSettings is an extension of Google PubSub's SubscriptionConfig that
// enables further configuration for a subscription of a topic
type SubscriptionSettings struct {
	EnableMessageOrdering bool
	RetainAckedMessages   bool
}

// TopicSettings is an extension of Google PubSub's TopicConfig that
// enables further configuration for a topic
type TopicSettings struct {
	RetentionDurationInDays int
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
func (ps *PubSub) CreateSubscriptions(tid string, sids map[string]string, ss SubscriptionSettings) error {
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
	// let's create subscriptions with optional filters for the topic
	for sid, flt := range sids {
		// check if the subscription exists or not
		s := ps.client.Subscription(sid)
		exists, err := s.Exists(ctx)
		if err != nil {
			return err
		}
		// if the sub already exists, skip it
		if exists {
			continue
		}
		// let's create a subscription. first, gather the configurations
		cfg := pubsub.SubscriptionConfig{
			EnableMessageOrdering: ss.EnableMessageOrdering,
			RetainAckedMessages:   ss.RetainAckedMessages,
			Topic:                 topic,
		}
		// only if the filter is provided in the map[string]string to include the Filter config
		if flt != "" {
			cfg.Filter = flt
		}
		// creating a subscription
		_, err = ps.client.CreateSubscription(ctx, sid, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create a Topic in Google PubSub if not exist
func (ps *PubSub) CreateTopic(tid string, ts TopicSettings) error {
	ctx := context.Background()

	topic := ps.client.Topic(tid)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	rdd := ts.RetentionDurationInDays
	if rdd > 7 {
		rdd = 7 // max to 7 days
	}
	cfg := pubsub.TopicConfig{
		RetentionDuration: time.Duration(rdd) * 24 * time.Hour,
	}
	_, err = ps.client.CreateTopicWithConfig(ctx, tid, &cfg)
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

	// if isLocal is true or is settingsPath is empty, then use credentials
	useCreds := isLocal || settingsPath != ""
	if useCreds {
		return pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(settingsPath))
	}

	return pubsub.NewClient(ctx, projectID)
}
