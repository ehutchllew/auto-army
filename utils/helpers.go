package utils

import (
	"fmt"
	"strconv"
)

func SafeConvertBool(val any) bool {
	if val == nil {
		return false // Return zero value of bool
	}

	return val.(bool)
}

func SafeConvertString(val any) string {
	if val == nil {
		return "" // Return zero value of string
	}

	return val.(string)
}

func SafeConvertUint8(val any) (uint8, error) {
	if val == nil {
		return 0, nil // Return zero value of uint8
	}
	switch v := val.(type) {
	case float64:
		return uint8(v), nil
	case int:
		return uint8(v), nil
	case int64:
		return uint8(v), nil
	case uint:
		return uint8(v), nil
	case uint64:
		return uint8(v), nil
	case float32:
		return uint8(v), nil
	case string:
		// If it's a string, try to parse it as a number
		parsed, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return 0, fmt.Errorf("failed to parse string as uint8: %v", err)
		}
		return uint8(parsed), nil
	default:
		return 0, fmt.Errorf("unsupported type for uint8 conversion: %T", val)
	}
}
