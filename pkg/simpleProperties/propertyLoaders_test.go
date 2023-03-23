package simpleProperties

import (
	"container/list"
	"reflect"
	"testing"
)

func TestBootPropertyLoader(t *testing.T) {
	type args struct {
		path string
	}
	bootProperties := make(map[string]string)
	bootProperties["application.name"] = "simpleProperties.app"
	bootProperties["boot.something"] = "123"
	tests := []struct {
		name string
		args args
		want Properties
	}{
		{
			name: "Test bootstrap loader",
			args: args{"testdata/resource/bootstrap"},
			want: Properties{
				bootProperties,
				make(map[string]string),
				make(map[string]string),
				make(map[string]*list.List),
				nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testProperties = &Properties{
				make(map[string]string),
				make(map[string]string),
				make(map[string]string),
				make(map[string]*list.List),
				nil}

			f := BootPropertyLoader(tt.args.path)
			f(testProperties)
			if !reflect.DeepEqual(*testProperties, tt.want) {
				t.Errorf("BootPropertyLoader() = %v, want %v", *testProperties, tt.want)
				t.Errorf("%v", *testProperties)
				t.Errorf("%v", tt.want)
			}
		})
	}
}
