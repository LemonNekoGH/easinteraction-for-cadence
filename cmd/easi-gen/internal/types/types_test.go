package types

import (
	"github.com/onflow/cadence/runtime/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStruct_Kind(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		s := &Struct{}
		assert.Equal(t, common.CompositeKindStructure, s.Kind())
	})

	t.Run("as super type", func(t *testing.T) {
		var s CompositeType
		s = &Struct{}
		assert.Equal(t, common.CompositeKindStructure, s.Kind())
	})
}

func TestResource_Kind(t *testing.T) {
	s := &Resource{}
	assert.Equal(t, common.CompositeKindResource, s.Kind())
}

func TestContract_Kind(t *testing.T) {
	s := &Contract{}
	assert.Equal(t, common.CompositeKindContract, s.Kind())
}

func TestEvent_Kind(t *testing.T) {
	s := &Event{}
	assert.Equal(t, common.CompositeKindEvent, s.Kind())
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
