package main

import (
	"fmt"
	"strings"
)

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}

func isNumber(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

func isString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

func isRune(v interface{}) bool {
	_, ok := v.([]rune)
	return ok
}
func isBool(v interface{}) bool {
	_, ok := v.(bool)
	return ok
}

func isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	if isNumber(left) && isNumber(right) {
		return left.(float64) == right.(float64)
	}
	if isString(left) && isString(right) {
		return left.(string) == right.(string)
	}
	if isRune(left) && isRune(right) {
		return string(left.([]rune)) == string(right.([]rune))
	}

	// isBool
	if isBool(left) && isBool(right) {
		return left.(bool) == right.(bool)
	}
	return false
}

func stringfy(obj any) string {
	switch v := obj.(type) {
	case float64:
		s := fmt.Sprintf("%f", v)

		// Trim trailing zeros
		s = strings.TrimRight(s, "0")

		// Trim the decimal point if it's the last character
		s = strings.TrimRight(s, ".")
		return s

	case string:
		return fmt.Sprintf("%s", v)
	case []rune:
		return fmt.Sprintf("%s", string(v))
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("%v", obj)
	}
}
