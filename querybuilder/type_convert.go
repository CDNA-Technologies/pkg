package querybuilder

import (
	"fmt"
	"strconv"
	"time"
)

const (
	DATE_ISO_8601      = "2006-01-02"
	TIME_ISO_8601      = "15:04:05"
	DATE_TIME_ISO_8601 = "2006-01-02T15:04:05"
)

//Error Checker function
func checkError(val interface{}, err interface{}) interface{} {
	if err == nil {
		return val
	}
	return nil
}

// String
func toString(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return ""
	}
}

// Double
func toDouble(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		return checkError(f, err)
	case float64:
		return v
	default:
		return 0
	}
}

// Integer
func toInteger(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		i, err := strconv.Atoi(v)
		return checkError(i, err)
	case float64:
		return int(v)
	case int:
		return v
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// Boolean
func toBoolean(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		b, err := strconv.ParseBool(v)
		return checkError(b, err)
	case float64:
		n := int(v)
		if n == 1 {
			return true
		}
		return false
	case bool:
		return v
	default:
		return false
	}
}

// Date
func toDate(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		t, err := time.Parse(DATE_ISO_8601, v)
		return checkError(t, err)
	default:
		return time.Time{}
	}
}

// Time
func toTime(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		t, err := time.Parse(TIME_ISO_8601, v)
		return checkError(t, err)
	default:
		return time.Time{}
	}
}

// DateTime
func toDateTime(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		t, err := time.Parse(DATE_TIME_ISO_8601, v)
		return checkError(t, err)
	default:
		return time.Time{}
	}
}
