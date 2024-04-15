package pb

import (
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// PubSubOptions provides a way to configure the PubSub client with various
// options such as the project ID, client options, and publish and receive settings.
type PubSubOptions struct {
	AutoOriginatedAt bool
	ProjectID        string
	ClientOptions    []option.ClientOption
	PublishSettings  pubsub.PublishSettings
	ReceiveSettings  pubsub.ReceiveSettings
}

// Options returns a new PubSubOptions struct with the provided project ID and
// client options. The AutoOriginatedAt field is set to true by default.
// Any ClientOptions provided are passed directly through to the underlying
// pubsub.Client when calling the NewPubSub function.
func Options(pID string, opts ...option.ClientOption) *PubSubOptions {
	return &PubSubOptions{
		AutoOriginatedAt: true,
		ProjectID:        pID,
		ClientOptions:    opts,
	}
}

// SetAutoOriginatedAt sets the AutoOriginatedAt field on the PubSubOptions struct
// to the provided value and returns the modified PubSubOptions struct. If true,
// the OriginatedAt attribute will be set to the current time if it is not already
// set for all messages published.
func (o *PubSubOptions) SetAutoOriginatedAt(auto bool) *PubSubOptions {
	o.AutoOriginatedAt = auto
	return o
}

// SetProjectID sets the ProjectID field on the PubSubOptions struct to the provided
// value and returns the modified PubSubOptions struct.
func (o *PubSubOptions) SetProjectID(pID string) *PubSubOptions {
	o.ProjectID = pID
	return o
}

// SetClientOptions sets the ClientOptions field on the PubSubOptions struct to the
// provided options and returns the modified PubSubOptions struct.
func (o *PubSubOptions) SetClientOptions(opts ...option.ClientOption) *PubSubOptions {
	if o.ClientOptions == nil {
		o.ClientOptions = []option.ClientOption{}
	}

	o.ClientOptions = append(o.ClientOptions, opts...)

	return o
}

// SetPublishSettings sets the PublishSettings field on the PubSubOptions struct to
// the provided settings and returns the modified PubSubOptions struct.
func (o *PubSubOptions) SetPublishSettings(s pubsub.PublishSettings) *PubSubOptions {
	o.PublishSettings = s
	return o
}

// SetReceiveSettings sets the ReceiveSettings field on the PubSubOptions struct to
// the provided settings and returns the modified PubSubOptions struct.
func (o *PubSubOptions) SetReceiveSettings(s pubsub.ReceiveSettings) *PubSubOptions {
	o.ReceiveSettings = s
	return o
}
