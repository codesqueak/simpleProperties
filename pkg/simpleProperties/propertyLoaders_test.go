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
	bootProperties["float"] = "123.45"
	bootProperties["int"] = "777"
	bootProperties["boolean"] = "true"
	bootProperties["boostrap.json"] = "test_json"
	bootProperties["boostrap.yaml"] = "test_yaml"
	tests := []struct {
		name string
		args args
		want Properties
	}{
		{
			name: "Test bootstrap loader",
			args: args{"testdata/resources/bootstrap"},
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
			var nameValue = testProperties.GetBootProperty("application.name")
			if nameValue != "simpleProperties.app" {
				t.Errorf("BootPropertyLoader() application.name = %v, want %v", nameValue, "application.name")
			}
		})
	}
}

func TestGlobalPropertyLoader(t *testing.T) {
	globalProperties := make(map[string]string)
	globalProperties["json1"] = "application.json"
	globalProperties["profile"] = "dev, debug"
	globalProperties["properties1"] = "application.properties"
	globalProperties["yaml1"] = "application.yaml"
	//
	expressionProperties := make(map[string]string)
	expressionProperties["expression"] = "An expression ${xyzzy}"
	//
	l := list.New()
	l.PushBack(&exprParts{"${xyzzy}", "xyzzy", ""})
	expressionList := make(map[string]*list.List)
	expressionList["expression"] = l

	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want Properties
	}{
		{
			name: "Test global property loader",
			args: args{"testdata/resources/application"},
			want: Properties{
				make(map[string]string),
				globalProperties,
				expressionProperties,
				expressionList,
				nil},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := EmptyProperties()
			f := GlobalPropertyLoader(tt.args.path)
			f(p)
			if !reflect.DeepEqual(*p, tt.want) {
				t.Errorf("GlobalPropertyLoader() = %v, want %v", *p, tt.want)
			}
		})
	}
}

func TestProfilePropertyLoader(t *testing.T) {
	profileProperties := make(map[string]string)
	profileProperties["profile"] = "test_profile"
	profileProperties["profile.property"] = "property profile value"

	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want Properties
	}{
		{
			name: "Test profile property loader",
			args: args{"testdata/resources/application"},
			want: Properties{
				make(map[string]string),
				profileProperties,
				make(map[string]string),
				make(map[string]*list.List),
				nil},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var baseProps = make(map[string]string)
			baseProps["profile"] = "test_profile"
			p := &Properties{make(map[string]string, 32),
				baseProps,
				make(map[string]string, 32),
				make(map[string]*list.List, 32),
				nil,
			}
			f := ProfilePropertyLoader(tt.args.path)
			f(p)
			if !reflect.DeepEqual(*p, tt.want) {
				t.Errorf("ProfilePropertyLoader() = %v, want %v", *p, tt.want)
			}
			var nameValue = p.GetProperty("profile.property")
			if nameValue != "property profile value" {
				t.Errorf("BootPropertyLoader() profile.property = %v, want %v", nameValue, "property profile value")
			}
		})
	}
}

func TestOSEnvironmentLoader(t *testing.T) {
	t.Run("Test OS environment property loader", func(t *testing.T) {
		p := EmptyProperties()
		f := LoadOSEnvironment()
		f(p)
		if len(p.keyValueMap) == 0 {
			t.Errorf("OSEnvironmentLoader() is empty - didn't expect that!")
		}
	})
}

func Test_extractExpressions(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want *list.List
	}{
		{"Expression extract empty",
			args{"hello unknown"},
			list.New(),
		},
		{"Expression extract only name",
			args{"hello ${person}"},
			func() *list.List {
				l := list.New()
				l.PushBack(&exprParts{"${person}", "person", ""})
				return l
			}(),
		},
		{"Expression extract with default",
			args{"hello ${person:unknown}"},
			func() *list.List {
				l := list.New()
				l.PushBack(&exprParts{"${person:unknown}", "person", "unknown"})
				return l
			}(),
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractExpressions(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractExpressions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setKV(t *testing.T) {
	type args struct {
		p     *Properties
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"set kv no whitespace",
			args{EmptyProperties(), "key", "value"},
		},
		{
			"set kv with spaces",
			args{EmptyProperties(), "  key  ", "  value  "},
		},
		{
			"set kv with tabs",
			args{EmptyProperties(), " \t key ", "  value\t  "},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setKV(tt.args.p, tt.args.key, tt.args.value)
			if tt.args.p.GetProperty("key") != "value" {
				t.Errorf("setKV() = %v, want %v", tt.args.p.GetProperty("key"), "value")
			}
		})
	}
}
