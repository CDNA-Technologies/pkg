package operator

import (
	"reflect"
	"strings"
)

func init() {
	AddOperator(EndsWith)
}

// EndsWith operator check if string ends with value of value param
var EndsWith = &Operator{
	Name: "ends_with",
	Evaluate: func(input, value interface{}) bool {
		rv := reflect.ValueOf(value)
		if rv.Kind() != reflect.String {
			return false
		}

		return strings.HasSuffix(input.(string), value.(string))
	},
}
