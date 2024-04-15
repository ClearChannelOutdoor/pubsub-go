package pb

import (
	"encoding/json"
	"reflect"
	"testing"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func TestPubSubOptions_Options(t *testing.T) {
	type args struct {
		pID  string
		opts []option.ClientOption
	}
	tests := []struct {
		name string
		args args
		want *PubSubOptions
	}{
		{
			"should create args with projectID",
			args{
				pID: "test-project",
			},
			&PubSubOptions{
				AutoOriginatedAt: true,
				ProjectID:        "test-project",
			},
		},
		{
			"should create args with projectID and client options",
			args{
				pID: "test-project",
				opts: []option.ClientOption{
					option.WithCredentialsFile("test.json"),
				},
			},
			&PubSubOptions{
				AutoOriginatedAt: true,
				ProjectID:        "test-project",
				ClientOptions: []option.ClientOption{
					option.WithCredentialsFile("test.json"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Options(tt.args.pID, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				wntJ, _ := json.MarshalIndent(tt.want, "", "  ")
				gotJ, _ := json.MarshalIndent(got, "", "  ")
				t.Errorf("Options():\n %s, \nwant:\n %s", gotJ, wntJ)
			}
		})
	}
}

func TestPubSubOptions_SetAutoOriginatedAt(t *testing.T) {
	type args struct {
		auto bool
	}
	tests := []struct {
		name string
		args args
		want *PubSubOptions
	}{
		{
			"should set auto originated at to true",
			args{
				auto: true,
			},
			&PubSubOptions{
				AutoOriginatedAt: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &PubSubOptions{}
			if got := o.SetAutoOriginatedAt(tt.args.auto); !reflect.DeepEqual(got, tt.want) {
				wntJ, _ := json.MarshalIndent(tt.want, "", "  ")
				gotJ, _ := json.MarshalIndent(got, "", "  ")
				t.Errorf("SetAutoOriginatedAt():\n %s, \nwant:\n %s", gotJ, wntJ)
			}
		})
	}
}

func TestPubSubOptions_SetProjectID(t *testing.T) {
	type args struct {
		pID string
	}
	tests := []struct {
		name string
		args args
		want *PubSubOptions
	}{
		{
			"should set projectID to specified string",
			args{
				pID: "test",
			},
			&PubSubOptions{
				ProjectID: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &PubSubOptions{}
			if got := o.SetProjectID(tt.args.pID); !reflect.DeepEqual(got, tt.want) {
				wntJ, _ := json.MarshalIndent(tt.want, "", "  ")
				gotJ, _ := json.MarshalIndent(got, "", "  ")
				t.Errorf("SetProjectID():\n %s, \nwant:\n %s", gotJ, wntJ)
			}
		})
	}
}

func TestPubSubOptions_SetClientOptions(t *testing.T) {
	type args struct {
		opts []option.ClientOption
	}
	tests := []struct {
		name string
		args args
		want *PubSubOptions
	}{
		{
			"should properly set client options",
			args{
				opts: []option.ClientOption{
					option.WithCredentialsFile("test.json"),
				},
			},
			&PubSubOptions{
				ClientOptions: []option.ClientOption{
					option.WithCredentialsFile("test.json"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &PubSubOptions{}
			if got := o.SetClientOptions(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				wntJ, _ := json.MarshalIndent(tt.want, "", "  ")
				gotJ, _ := json.MarshalIndent(got, "", "  ")
				t.Errorf("SetClientOptions():\n %s, \nwant:\n %s", gotJ, wntJ)
			}
		})
	}
}

func TestPubSubOptions_SetPublishSettings(t *testing.T) {
	type args struct {
		s pubsub.PublishSettings
	}
	tests := []struct {
		name string
		args args
		want *PubSubOptions
	}{
		{
			"should properly set publish settings",
			args{
				pubsub.PublishSettings{
					ByteThreshold: 1024,
				},
			},
			&PubSubOptions{
				PublishSettings: pubsub.PublishSettings{
					ByteThreshold: 1024,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &PubSubOptions{}
			if got := o.SetPublishSettings(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				wntJ, _ := json.MarshalIndent(tt.want, "", "  ")
				gotJ, _ := json.MarshalIndent(got, "", "  ")
				t.Errorf("SetPublishSettings():\n %s, \nwant:\n %s", gotJ, wntJ)
			}
		})
	}
}

func TestPubSubOptions_SetReceiveSettings(t *testing.T) {
	type args struct {
		s pubsub.ReceiveSettings
	}
	tests := []struct {
		name string
		args args
		want *PubSubOptions
	}{
		{
			"should properly set receive settings",
			args{
				pubsub.ReceiveSettings{
					MaxOutstandingMessages: 10,
				},
			},
			&PubSubOptions{
				ReceiveSettings: pubsub.ReceiveSettings{
					MaxOutstandingMessages: 10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &PubSubOptions{}
			if got := o.SetReceiveSettings(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				wntJ, _ := json.MarshalIndent(tt.want, "", "  ")
				gotJ, _ := json.MarshalIndent(got, "", "  ")
				t.Errorf("SetReceiveSettings():\n %s, \nwant:\n %s", gotJ, wntJ)
			}
		})
	}
}
