package pb

import (
	"reflect"
	"testing"
)

func Test_mergeMaps(t *testing.T) {
	type args struct {
		attrs []map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			"should return empty map when no maps are provided",
			args{},
			map[string]string{},
		},
		{
			"should return a map that is the same when only 1 map is provided",
			args{
				[]map[string]string{
					{"key": "value"},
				},
			},
			map[string]string{
				"key": "value",
			},
		},
		{
			"should return a map that combines all maps that are provided",
			args{
				[]map[string]string{
					{"key": "value"},
					{"key2": "value2"},
					{
						"key3": "value3",
						"key4": "value4",
					},
				},
			},
			map[string]string{
				"key":  "value",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeMaps(tt.args.attrs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}
