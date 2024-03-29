package pb

import (
	"encoding/json"
	"google.golang.org/api/option"
	"reflect"
	"testing"
)

func TestOptions(t *testing.T) {
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
