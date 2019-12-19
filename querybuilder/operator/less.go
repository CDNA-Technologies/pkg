package operator

import "time"

func init() {
	AddOperator(Less)
}

var Less = &Operator{
	Name: "less",
	Evaluate: func(input, value interface{}) bool {
		if input == nil {
			return false
		}

		switch value.(type) {
		case float64:
			return input.(float64) < value.(float64)
		case int:
			return input.(int) < value.(int)
		case string:
			return len(input.(string)) < len(value.(string))
		case time.Time:
			i := input.(time.Time)
			v := value.(time.Time)
			return i.Before(v)
		default:
			return false
		}
	},
}
