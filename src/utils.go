package src

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// ConvertTypeToProto converts go data type to protobuf data type.
func ConvertTypeToProto(t string) string {
	switch t {
	case "int", "int64":
		return "int64"
	case "uint", "uint64":
		return "uint64"
	case "int32":
		return "int32"
	case "uint32":
		return "uint32"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "bool":
		return "bool"
	case "string":
		return "string"
	default:
		return t
	}
}

// ConvertTypeToGraphQL converts go data type to GraphQL data type.
func ConvertTypeToGraphQL(t string) string {
	switch t {
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		return "Int"
	case "float32", "float64":
		return "Float"
	case "bool":
		return "Boolen"
	case "string":
		return "String"
	default:
		return t
	}
}

// ReplaceFileExt replaces ext of filePath.
func ReplaceFileExt(filePath, ext string) string {
	rawExt := filepath.Ext(filePath)
	if rawExt == "" {
		return fmt.Sprintf("%s%s", filePath, ext)
	}
	return strings.Replace(filePath, rawExt, ext, 1)
}

// GetStructTag finds json tag from given string, and returns key and omitempty or not.
func GetStructTag(s string) (key string, omitempty bool) {
	const empty = "omitempty"
	// json or bson
	re := regexp.MustCompile(`.*son:"(.+)"`)
	slice := re.FindStringSubmatch(s)
	if len(slice) != 2 {
		return
	}

	value := slice[1]
	if strings.Contains(value, empty) {
		keySlice := strings.Split(value, ",")
		if len(keySlice) == 1 {
			return "", true
		}
		for i := range keySlice {
			if keySlice[i] != empty {
				return keySlice[i], true
			}
		}
	}
	return value, false
}
