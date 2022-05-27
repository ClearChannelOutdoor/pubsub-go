package pubsub_go

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"google.golang.org/api/option"
)

// PubSub TODO: Description
type PubSub struct {
	client   *pubsub.Client
	settings Settings
}

type Settings struct {
	publish pubsub.PublishSettings
	receive pubsub.ReceiveSettings
}

// NewPubSub TODO: Description
func NewPubSub(projectID string, settings Settings, isLocal bool, settingsPath string) (*PubSub, error) {
	client, err := newClient(projectID, isLocal, settingsPath)
	if err != nil {
		return nil, err
	}

	return &PubSub{
		client:   client,
		settings: settings,
	}, nil
}

// Publish TODO: Description
func (ps *PubSub) Publish(m Message) error {
	topic := ps.client.Topic(m.Topic)
	topic.PublishSettings = ps.settings.publish

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

func newClient(projectID string, isLocal bool, settingsPath string) (*pubsub.Client, error) {
	ctx := context.Background()

	if isLocal {
		return pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(settingsPath))
	}

	return pubsub.NewClient(ctx, projectID)
}
