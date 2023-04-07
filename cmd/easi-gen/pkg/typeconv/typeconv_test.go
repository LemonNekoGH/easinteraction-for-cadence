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
				ok, goName, _ := MaybeMapType(typeName, nil)
				expected := fmt.Sprintf("map[%s]%s", v, v2)
				if !ok || goName != expected {
					t.Errorf("expected: %v, %s, got: %v, %s", true, expected, ok, goName)
				}
			})
		}
	}

	t.Run("nested map: {String:{String:{String:{String:{String:{String:String}}}}}}", func(t *testing.T) {
		ok, goName, _ := MaybeMapType("{String:{String:{String:{String:{String:{String:String}}}}}}", nil)
		expected := "map[string]map[string]map[string]map[string]map[string]map[string]string"
		if !ok || goName != expected {
			t.Errorf("expected: %v, %s, got: %v, %s", true, expected, ok, goName)
		}
	})
}

func TestMaybeArrayType(t *testing.T) {
	for k, v := range typeMap {
		typeName := fmt.Sprintf("[%s]", k)
		t.Run(fmt.Sprintf("simple array: %s", typeName), func(t *testing.T) {
			ok, goName, _ := MaybeArrayType(typeName, nil)
			expected := fmt.Sprintf("[]%s", v)
			if !ok || goName != expected {
				t.Errorf("expected: %v, %s, got: %v, %s", true, expected, ok, goName)
			}
		})
	}

	t.Run("nested array: [[[String]]]", func(t *testing.T) {
		ok, goName, _ := MaybeArrayType("[[[String]]]", nil)
		expected := "[][][]string"
		if !ok || goName != expected {
			t.Errorf("expected: %v, %s, got: %v, %s", true, expected, ok, goName)
		}
	})

	t.Run("fixed size array", func(t *testing.T) {
		ok, goName, _ := MaybeArrayType("[String; 4]", nil)
		expected := "[4]string"
		if !ok || goName != expected {
			t.Errorf("expected: %v, %s, got: %v, %s", true, expected, ok, goName)
		}
	})
}
