package pb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
)

type gcpPubSubProvider interface {
	Close() error
	CreateSubscription(context.Context, string, pubsub.SubscriptionConfig) (*pubsub.Subscription, error)
	CreateTopic(context.Context, string) (*pubsub.Topic, error)
	CreateTopicWithConfig(context.Context, string, *pubsub.TopicConfig) (*pubsub.Topic, error)
	Subscription(string) *pubsub.Subscription
	Topic(string) *pubsub.Topic
}

type PubSub struct {
	clnt gcpPubSubProvider
	ctx  context.Context
	opts *PubSubOptions
}

func (p *PubSub) ensurePublishSettings() {
	if p.opts.PublishSettings == (pubsub.PublishSettings{}) {
		p.opts.PublishSettings = pubsub.DefaultPublishSettings
	}
}

func (p *PubSub) ensureReceiveSettings() {
	if p.opts.ReceiveSettings == (pubsub.ReceiveSettings{}) {
		p.opts.ReceiveSettings = pubsub.DefaultReceiveSettings
	}
}

func (p *PubSub) Close() error {
	return p.clnt.Close()
}

func (p *PubSub) CreateSubscription(id string, sid string, fltr string, cfg ...pubsub.SubscriptionConfig) error {
	// ensure we have a subscription config
	ss := pubsub.SubscriptionConfig{}
	if len(cfg) > 0 {
		ss = cfg[0]
	}

	// set topic for the subscription
	t := p.clnt.Topic(id)
	ss.Topic = t

	// check to see if the requested subscription already exists
	exists, err := p.clnt.Subscription(sid).Exists(p.ctx)
	if err != nil {
		return err
	}

	// create the subscription if it does not exist
	if !exists {
		// set the filter if provided
		if fltr != "" {
			ss.Filter = fltr
		}

		// create the subscription
		if _, err := p.clnt.CreateSubscription(p.ctx, sid, ss); err != nil {
			return err
		}
	}

	return nil
}

func (p *PubSub) CreateSubscriptions(id string, sids map[string]string, cfg ...pubsub.SubscriptionConfig) error {
	// create each subscription
	for sid, f := range sids {
		if err := p.CreateSubscription(id, sid, f, cfg...); err != nil {
			return err
		}
	}

	return nil
}

func (p *PubSub) CreateTopic(id string, cfg ...pubsub.TopicConfig) error {
	// check to see if the requested topic already exists
	exists, err := p.clnt.Topic(id).Exists(p.ctx)
	if err != nil {
		return err
	}

	// create the topic with or without configuration
	ct := func() error {
		if len(cfg) > 0 {
			if _, err := p.clnt.CreateTopicWithConfig(p.ctx, id, &cfg[0]); err != nil {
				return err
			}
		}

		if _, err := p.clnt.CreateTopic(p.ctx, id); err != nil {
			return err
		}

		return nil
	}

	// create the topic if it does not exist
	if !exists {
		if err := ct(); err != nil {
			return err
		}
	}

	// topic exists
	return nil
}

func (p *PubSub) Publish(id string, d any, attrs ...map[string]string) error {
	t := p.clnt.Topic(id)

	// apply PublishSettings
	p.ensurePublishSettings()
	t.PublishSettings = p.opts.PublishSettings

	// marshal provided data as JSON if needed
	var dta []byte
	if _, ok := d.([]byte); !ok {
		data, err := json.Marshal(d)
		if err != nil {
			return err
		}

		dta = data
	}

	// ensure data to published is set
	if dta == nil {
		dta = d.([]byte)
	}

	// apply OriginatedAt attribute
	mgd := mergeMaps(attrs...)

	// set OriginatedAt attribute if not set and AutoOriginatedAt is true
	if _, ok := mgd["OriginatedAt"]; p.opts.AutoOriginatedAt && !ok {
		mgd["OriginatedAt"] = fmt.Sprintf("%v", time.Now().Unix())
	}

	// publish the message
	res := t.Publish(p.ctx, &pubsub.Message{
		Data:       dta,
		Attributes: mgd,
	})

	// get the result to ensure message was published
	if _, err := res.Get(p.ctx); err != nil {
		return err
	}

	return nil
}

func (p *PubSub) Receive(id string, mc chan<- *pubsub.Message) error {
	sub := p.clnt.Subscription(id)
	p.ensureReceiveSettings()
	sub.ReceiveSettings = p.opts.ReceiveSettings

	return sub.Receive(p.ctx, func(ctx context.Context, m *pubsub.Message) {
		mc <- m
	})
}

func NewPubSub(ctx context.Context, opts *PubSubOptions) (*PubSub, error) {
	clnt, err := pubsub.NewClient(ctx, opts.ProjectID, opts.ClientOptions...)
	if err != nil {
		return nil, err
	}

	return &PubSub{
		clnt: clnt,
		ctx:  ctx,
		opts: opts,
	}, nil
}
