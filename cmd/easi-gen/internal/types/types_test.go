package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStruct_IsStructOrIsResource(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		s := &Struct{}
		if !s.IsStruct() || s.IsResource() {
			t.Fail()
		}
		assert.True(t, s.IsStruct())
		assert.False(t, s.IsResource())
	})

	t.Run("as super type", func(t *testing.T) {
		var s CompositeType
		s = &Struct{}
		if !s.IsStruct() || s.IsResource() {
			t.Fail()
		}
		assert.True(t, s.IsStruct())
		assert.False(t, s.IsResource())
	})
}

func TestResource_IsStructOrIsResource(t *testing.T) {
	s := &Resource{}
	if s.IsStruct() || !s.IsResource() {
		t.Fail()
	}
}

func TestContract_IsStructOrIsResource(t *testing.T) {
	s := &Contract{}
	if s.IsStruct() || s.IsResource() {
		t.Fail()
	}
}

func TestContract_FlattenSubTypes(t *testing.T) {
	c := &Contract{}
	s := &Struct{}
	r := &Resource{}

	s.SetSubTypes([]CompositeType{r})
	c.SetSubTypes([]CompositeType{s})
	c.FlattenSubTypes()

	assert.ElementsMatch(t, []CompositeType{r, s}, c.GetSubTypes())
}

func TestFunction_AuthorizerCount(t *testing.T) {
	f1 := Function{}
	assert.Equal(t, 0, f1.AuthorizerCount())
	f1.Params = []FunctionParam{{Type: "AuthAccount"}}
	assert.Equal(t, 1, f1.AuthorizerCount())
}
