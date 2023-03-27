package gen

import "github.com/onflow/cadence"

var (
	typeMap = map[string]string{
		cadence.TheVoidType.ID():           "any",
		cadence.TheBoolType.ID():           "bool",
		cadence.TheAddressType.ID():        "string",
		cadence.TheStringType.ID():         "string",
		cadence.TheIntType.ID():            "*big.Int",
		cadence.TheInt8Type.ID():           "int8",
		cadence.TheInt16Type.ID():          "int16",
		cadence.TheInt32Type.ID():          "int32",
		cadence.TheInt64Type.ID():          "int64",
		cadence.TheUIntType.ID():           "*big.Int",
		cadence.TheUInt8Type.ID():          "uint8",
		cadence.TheUInt16Type.ID():         "uint16",
		cadence.TheUInt32Type.ID():         "uint32",
		cadence.TheUInt64Type.ID():         "uint64",
		cadence.TheWord8Type.ID():          "uint8",
		cadence.TheWord32Type.ID():         "uint32",
		cadence.TheWord64Type.ID():         "uint64",
		cadence.TheUFix64Type.ID():         "uint64",
		cadence.TheFix64Type.ID():          "int64",
		cadence.TheInt128Type.ID():         "*big.Int",
		cadence.TheUInt128Type.ID():        "*big.Int",
		cadence.TheInt256Type.ID():         "*big.Int",
		cadence.TheUInt256Type.ID():        "*big.Int",
		cadence.TheCapabilityPathType.ID(): "string",
		cadence.TheAnyStructType.ID():      "any",
		cadence.TheAnyResourceType.ID():    "any",
	}
)
