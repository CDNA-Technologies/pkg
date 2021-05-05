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

// String
func toString(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%f", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	default:
		return "", nil
	}
}

// Double
func toDouble(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	default:
		return 0, nil
	}
}

// Integer
func toInteger(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case string:
		return strconv.Atoi(v)
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, nil
	}
}

// Boolean
func toBoolean(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseBool(v)
	case float64:
		n := int(v)
		if n == 1 {
			return true, nil
		}
		return false, nil
	case bool:
		return v, nil
	default:
		return false, nil
	}
}

// Date
func toDate(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case string:
		return time.Parse(DATE_ISO_8601, v)
	default:
		return time.Time{}, nil
	}
}

// Time
func toTime(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case string:
		return time.Parse(TIME_ISO_8601, v)
	default:
		return time.Time{}, nil
	}
}

// DateTime
func toDateTime(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case string:
		return time.Parse(DATE_TIME_ISO_8601, v)
	default:
		return time.Time{}, nil
	}
}
