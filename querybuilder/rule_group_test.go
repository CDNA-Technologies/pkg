package querybuilder

import (
	"encoding/json"
	"testing"
)

func parseJson(data string) map[string]interface{} {
	var dataset map[string]interface{}
	if err := json.Unmarshal([]byte(data), &dataset); err != nil {
		return nil
	}

	return dataset
}

var ruleset = []string{
	`{
  "condition": "AND",
  "rules": [
		{"id": "float_equal","field": "float_equal","type": "double","input": "text","operator": "equal","value": 1.2},
		{"id": "int_equal","field": "int_equal","type": "integer","input": "text","operator": "equal","value": 5},
		{"id": "int_greater","field": "int_greater","type": "integer","input": "text","operator": "greater","value": 2},
		{"id": "float_greater","field": "float_greater","type": "double","input": "text","operator": "greater","value": 7.5}
	]}`,

	`{
  "condition": "OR",
  "rules": [
		{"id": "float_equal","field": "float_equal","type": "double","input": "text","operator": "equal","value": 1.2},
		{"id": "int_greater","field": "int_greater","type": "integer","input": "text","operator": "greater","value": 2},
		{"id": "int_equal","field": "int_equal","type": "integer","input": "text","operator": "equal","value": 5}
	]}`,

	`{
  "condition": "AND",
  "rules": [
		{"id": "float_equal","field": "float_equal","type": "double","input": "text","operator": "equal","value": 1.2},
		{"id": "float_greater","field": "float_greater","type": "double","input": "text","operator": "greater","value": 7.5},
		{"condition": "OR","rules": [
			{"id": "int_greater","field": "int_greater","type": "integer","input": "text","operator": "greater","value": 2},
			{"id": "int_equal","field": "int_equal","type": "integer","input": "text","operator": "equal","value": 5}
		]}
	]}`,

	`{
	"condition": "AND",
	"rules": [
		  {"id": "float_greater_or_equal","field": "float_greater_or_equal","type": "double","input": "text","operator": "greater_or_equal","value": 1.2},
		  {"id": "datetime_greater","field": "datetime_greater","type": "datetime","input": "text","operator": "greater","value": "2021-01-01T04:21:21"},
		  {"condition": "OR","rules": [
			  {"id": "time_lesser","field": "time_lesser","type": "time","input": "text","operator": "less","value": "09:09:21"},
			  {"id": "int_equal","field": "int_equal","type": "integer","input": "text","operator": "equal","value": 5}
		  ]}
	  ]}`,

	`{
	"condition": "OR",
	"rules": [
			{"id": "int_lesser","field": "int_lesser","type": "integer","input": "text","operator": "less","value": 1},
			{"id": "string_contains ","field": "string_contains","type": "string","input": "text","operator": "contains","value": "test"},
			{"condition": "AND","rules": [
				{"id": "time_lesser","field": "time_lesser","type": "time","input": "text","operator": "less","value": "09:09:21"},
				{"id": "int_equal","field": "int_equal","type": "integer","input": "text","operator": "equal","value": 50}
			]}
		]}`,
}

func TestRuleGroupEvaluate(t *testing.T) {
	inputs := []struct {
		title   string
		rules   string
		dataset string
		want    bool
	}{
		{"(1) only AND", ruleset[0], `{"float_equal":  1.2, "int_equal": 5, "int_greater":  3,"float_greater": 7.7}`, true},
		{"(2) only AND", ruleset[0], `{"float_equal":  1.2, "int_equal": 4, "int_greater":  3,"float_greater": 7.7}`, false},
		{"(3) only OR", ruleset[1], `{"float_equal":  1.3, "int_equal": 4.5, "int_greater":  1, "float_greater": 7.7}`, false},
		{"(4) only OR", ruleset[1], `{"float_equal":  1.3, "int_equal": 4, "int_greater":  3, "float_greater": 7.7}`, true},
		{"(5) AND & OR", ruleset[2], `{"float_equal":  1.2, "int_equal": 5, "int_greater":  1,"float_greater": 7.7}`, true},
		{"(6) AND & OR", ruleset[2], `{"float_equal":  1.2, "int_equal": 5, "int_greater":  3,"float_greater": 7.9}`, true},
		{"(7) AND & OR", ruleset[3], `{"float_greater_or_equal":  "1.2", "datetime_greater": "2021-02-01T21:21:24", "time_lesser": "02:04:59" ,"int_equal": 10}`, true},
		{"(7) AND & OR", ruleset[3], `{"float_greater_or_equal":  "20a", "datetime_greater": "2020-02-01T21:21:24", "time_lesser": "02:04:59" ,"int_equal": 10}`, false},
		{"(8) OR & AND", ruleset[4], `{"int_lesser":  "2", "string_contains": "test_cricket", "time_lesser": "02:04:59" ,"int_equal": 10}`, true},
	}

	for _, input := range inputs {
		t.Run(input.title, func(t *testing.T) {
			rules := parseJson(input.rules)
			if rules == nil {
				t.Fatal("not parse json")
			}

			rg := RuleGroup{Condition: rules["condition"], Rules: rules["rules"]}

			if got, err := rg.Evaluate(parseJson(input.dataset)); got != input.want {
				t.Errorf("Caught error %s and returned %t", err, got)
			}
		})
	}
}

func BenchmarkRuleGroupEvaluate(b *testing.B) {
	rules := parseJson(ruleset[2])
	if rules == nil {
		b.Fatal("not parse json")
	}

	rg := RuleGroup{Condition: rules["condition"], Rules: rules["rules"]}
	dataset := parseJson(`{"float_equal":  1.2, "int_equal": 5, "int_greater":  1,"float_greater": 7.7}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rg.Evaluate(dataset)
	}
}
