package gen

import (
	_ "embed"
	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_contractFunction_commaCountAll(t *testing.T) {
	c := types.Function{}
	assert.Equal(t, -1, c.CommaCountAll())
}

func Test_contractFunction_commaCountCommon(t *testing.T) {
	c := types.Function{
		Params: []types.FunctionParam{{
			Type: "AuthAccount",
		}},
	}
	assert.Equal(t, -1, c.CommaCountCommon())
}

func Test_contractFunction_commaCountAuth(t *testing.T) {
	c := types.Function{
		Params: []types.FunctionParam{{
			Type: "AuthAccount",
		}},
	}
	assert.Equal(t, 0, c.CommaCountAuth())
}

func Test_contractFunction_genCadenceScript(t *testing.T) {
	t.Run("tx", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		fn := types.Function{
			OwnerTypeName: "AContract",
			Name:          "setSomething",
			GoName:        "SetSomething",
			Params: []types.FunctionParam{
				{
					Label:  "firstArg",
					Name:   "arg0",
					Type:   "AuthAccount",
					GoType: "",
				},
				{
					Label:  "_",
					Name:   "arg1",
					Type:   "UInt8",
					GoType: "",
				},
			},
		}

		generated, err := fn.GenCadenceScript()
		r.Empty(err)
		a.Equal(`import AContract from %s
transaction(arg1: UInt8) {
    prepare(arg0: AuthAccount) {
        AContract.setSomething(firstArg:arg0,arg1)
    }
}
`, generated)
	})

	t.Run("query", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		fn := types.Function{
			OwnerTypeName: "AContract",
			Name:          "setSomething",
			GoName:        "SetSomething",
			Params: []types.FunctionParam{
				{
					Label:  "firstArg",
					Name:   "arg",
					Type:   "Address",
					GoType: "",
				},
				{
					Label:  "_",
					Name:   "arg1",
					Type:   "UInt8",
					GoType: "",
				},
			},
			ReturnType: "String",
		}

		generated, err := fn.GenCadenceScript()
		r.Empty(err)
		a.Equal(`import AContract from %s
pub fun main(arg0: Address,arg1: UInt8): String{
    return AContract.setSomething(arg0,arg1)
}
`, generated)
	})
}
