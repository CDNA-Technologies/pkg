package operator

func init() {
	AddOperator(Greater)
}

// Greater
var Greater = &Operator{
	Name: "greater",
	Evaluate: func(input, value interface{}) bool {
		if input == nil {
			return false
		}

		switch value.(type) {
		case float64:
			return input.(float64) > value.(float64)
		case int:
			return input.(int) > value.(int)
		case string:
			return len(input.(string)) > len(value.(string))
		default:
			return false
		}
	},
}
