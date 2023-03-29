package typeconv

import (
	"fmt"
	"testing"
)

func Test_maybeMapType(t *testing.T) {
	for k, v := range typeMap {
		for k2, v2 := range typeMap {
			typeName := fmt.Sprintf("{%s:%s}", k, k2)
			t.Run(fmt.Sprintf("simple map: %s", typeName), func(t *testing.T) {
				ok, goName := MaybeMapType(typeName)
				expected := fmt.Sprintf("map[%s]%s", v, v2)
				if !ok || goName != expected {
					t.Errorf("expected: %v, %s, got: %v, %s", true, expected, ok, goName)
				}
			})
		}
	}

	t.Run("nested map: {String:{String:{String:{String:{String:{String:String}}}}}}", func(t *testing.T) {
		ok, goName := MaybeMapType("{String:{String:{String:{String:{String:{String:String}}}}}}")
		expected := "map[string]map[string]map[string]map[string]map[string]map[string]string"
		if !ok || goName != expected {
			t.Errorf("expected: %v, %s, got: %v, %s", true, expected, ok, goName)
		}
	})
}
