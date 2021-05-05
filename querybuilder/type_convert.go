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
	var err error = nil
	switch v := v.(type) {
	case string:
		return v, err
	case float64:
		return fmt.Sprintf("%f", v), err
	case bool:
		return fmt.Sprintf("%t", v), err
	default:
		return "", err
	}
}

// Double
func toDouble(v interface{}) (interface{}, error) {
	var err error = nil
	switch v := v.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, err
	default:
		return 0, err
	}
}

// Integer
func toInteger(v interface{}) (interface{}, error) {
	var err error = nil
	switch v := v.(type) {
	case string:
		return strconv.Atoi(v)
	case float64:
		return int(v), err
	case int:
		return v, err
	case bool:
		if v {
			return 1, err
		}
		return 0, err
	default:
		return 0, err
	}
}

// Boolean
func toBoolean(v interface{}) (interface{}, error) {
	var err error = nil
	switch v := v.(type) {
	case string:
		return strconv.ParseBool(v)
	case float64:
		n := int(v)
		if n == 1 {
			return true, err
		}
		return false, err
	case bool:
		return v, err
	default:
		return false, err
	}
}

// Date
func toDate(v interface{}) (interface{}, error) {
	var err error = nil
	switch v := v.(type) {
	case string:
		return time.Parse(DATE_ISO_8601, v)
	default:
		return time.Time{}, err
	}
}

// Time
func toTime(v interface{}) (interface{}, error) {
	var err error = nil
	switch v := v.(type) {
	case string:
		return time.Parse(TIME_ISO_8601, v)
	default:
		return time.Time{}, err
	}
}

// DateTime
func toDateTime(v interface{}) (interface{}, error) {
	var err error = nil
	switch v := v.(type) {
	case string:
		return time.Parse(DATE_TIME_ISO_8601, v)
	default:
		return time.Time{}, err
	}
}
