package main

import "fmt"

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
	return false
}

func stringfy(obj any) string {
	switch v := obj.(type) {
	case float64:
	    return fmt.Sprintf("%g", v) 
		// if v == float64(int(v)) {
		// 	return fmt.Sprintf("%.1f", v) // Ensures 1234.0 for whole numbers
		// } else {
		// 	return fmt.Sprintf("%g", v) // Keeps the precision for non-whole numbers
		// }
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
