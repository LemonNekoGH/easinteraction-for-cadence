package typeconv

import (
	"fmt"
	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/internal/types"
	"github.com/onflow/cadence"
	"github.com/samber/lo"
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

// MaybeOptionalType checks param is Cadence optional type or not
func MaybeOptionalType(cdcType string, foundTypes []types.CompositeType) (bool, string, string) {
	if !strings.HasSuffix(cdcType, "?") {
		return false, "", ""
	}
	goType, expanded := ByName(strings.TrimSuffix(cdcType, "?"), foundTypes)
	return true, goType, expanded + "?"
}

// MaybeCompositeType checks param is User's custom type
func MaybeCompositeType(cdcType string, foundTypes []types.CompositeType) (bool, string, string) {
	t, ok := lo.Find(foundTypes, func(it types.CompositeType) bool {
		return it.GetSimpleName() == cdcType
	})
	if !ok {
		return false, "", ""
	}
	return true, t.GetGoName(), t.GetName()
}

// MaybeMapType checks param is Cadence map or not.
func MaybeMapType(cdcType string, foundTypes []types.CompositeType) (bool, string, string) {
	if !strings.HasPrefix(cdcType, "{") || !strings.HasSuffix(cdcType, "}") {
		// not map
		return false, "", ""
	}
	typesStr := cdcType[1 : len(cdcType)-1]
	typesArr := strings.Split(typesStr, ":")
	keyType := typesArr[0]
	valueType := cdcType[len(keyType)+2 : len(cdcType)-1]
	keyGoType, keyExpanded := ByName(keyType, foundTypes)
	valueGoType, valueExpanded := ByName(valueType, foundTypes)
	return true, fmt.Sprintf("map[%s]%s", keyGoType, valueGoType), fmt.Sprintf("{%s:%s}", keyExpanded, valueExpanded)
}

// MaybeArrayType checks param is Cadence array or not.
func MaybeArrayType(cdcType string, foundTypes []types.CompositeType) (bool, string, string) {
	if !strings.HasPrefix(cdcType, "[") || !strings.HasSuffix(cdcType, "]") {
		// not array
		return false, "", ""
	}
	// check size
	typeAndSize := strings.Split(cdcType[1:len(cdcType)-1], ";")
	if len(typeAndSize) == 2 {
		size, err := strconv.Atoi(strings.TrimSpace(typeAndSize[1]))
		if err != nil {
			return false, "", ""
		}
		goType, expanded := ByName(typeAndSize[0], foundTypes)
		return true, fmt.Sprintf("[%d]%s", size, goType), fmt.Sprintf("[%s]", expanded)
	}
	goType, expanded := ByName(typeAndSize[0], foundTypes)
	return true, fmt.Sprintf("[]%s", goType), fmt.Sprintf("[%s]", expanded)
}

// ByName receives Cadence type name then returns corresponding Go type name.
func ByName(cdcType string, foundTypes []types.CompositeType) (string, string) {
	t := strings.TrimSpace(cdcType) // because map type name will be like {String: String}, there is a space character before value type
	if ok, t2, e := MaybeMapType(t, foundTypes); ok {
		return t2, e
	}
	if ok, t2, e := MaybeArrayType(cdcType, foundTypes); ok {
		return t2, e
	}
	if ok, t2, e := MaybeOptionalType(cdcType, foundTypes); ok {
		return t2, e
	}
	if ok, t2, e := MaybeCompositeType(cdcType, foundTypes); ok {
		return t2, e
	}

	return typeMap[t], cdcType
}
