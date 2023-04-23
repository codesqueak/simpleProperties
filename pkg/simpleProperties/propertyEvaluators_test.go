package simpleProperties

import (
	"container/list"
	"testing"
)

func TestBasicEvaluator(t *testing.T) {

	t.Run("Basic evaluator", func(t *testing.T) {

		var baseProps = make(map[string]string)
		baseProps["profile"] = "evaluator_profile"
		p := &Properties{make(map[string]string, 32),
			baseProps,
			make(map[string]string, 32),
			make(map[string]*list.List, 32),
			nil,
		}
		// load & evaluate
		loader := ProfilePropertyLoader("testdata/resources/application")
		loader(p)
		evaluator := BasicEvaluator()
		evaluator(p)
		//
		r := p.GetProperty("level1")
		want := "test value 1 and quark and default_v3 and default_v4 and xyzzy"
		if r != want {
			t.Errorf("BasicEvaluator() = %v, want %v", r, want)
		}
	})
}
