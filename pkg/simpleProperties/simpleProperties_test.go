package simpleProperties

import (
	"container/list"
	"reflect"
	"testing"
)

func TestNewProperties(t *testing.T) {
	tests := []struct {
		name string
		want Properties
	}{
		{name: "Test new properties",
			want: Properties{
				make(map[string]string),
				make(map[string]string),
				make(map[string]string),
				make(map[string]*list.List),
				nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := *EmptyProperties(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProperties() = %v, want %v", got, tt.want)
			}
		})
	}
}
