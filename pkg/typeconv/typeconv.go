package typeconv

import (
	"fmt"
	"github.com/onflow/cadence"
	"strconv"
	"strings"
)

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

// MaybeMapType checks param is Cadence map or not.
func MaybeMapType(cdcType string) (bool, string) {
	if cdcType[0] != '{' || cdcType[len(cdcType)-1] != '}' {
		// not map
		return false, ""
	}
	typesStr := cdcType[1 : len(cdcType)-1]
	types := strings.Split(typesStr, ":")
	keyType := types[0]
	valueType := cdcType[len(keyType)+2 : len(cdcType)-1]
	return true, fmt.Sprintf("map[%s]%s", ByName(keyType), ByName(valueType))
}

// MaybeArrayType checks param is Cadence array or not.
func MaybeArrayType(cdcType string) (bool, string) {
	if cdcType[0] != '[' || cdcType[len(cdcType)-1] != ']' {
		// not array
		return false, ""
	}
	// check size
	typeAndSize := strings.Split(cdcType[1:len(cdcType)-1], ";")
	if len(typeAndSize) == 2 {
		size, err := strconv.Atoi(strings.TrimSpace(typeAndSize[1]))
		if err != nil {
			return false, ""
		}
		return true, fmt.Sprintf("[%d]%s", size, ByName(typeAndSize[0]))
	}
	return true, fmt.Sprintf("[]%s", ByName(typeAndSize[0]))
}

// ByName receives Cadence type name then returns corresponding Go type name.
func ByName(cdcType string) string {
	t := strings.TrimSpace(cdcType) // because map type name will be like {String: String}, there is a space character before value type
	if ok, t2 := MaybeMapType(t); ok {
		return t2
	}
	if ok, t2 := MaybeArrayType(cdcType); ok {
		return t2
	}
	return typeMap[t]
}
