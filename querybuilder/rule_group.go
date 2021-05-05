package querybuilder

import "github.com/pkg/errors"

const (
	AND = "AND"
	OR  = "OR"
)

type Checker interface {
	Evaluate(dataset map[string]interface{}) (bool, error)
}

type RuleGroup struct {
	Condition interface{}
	Rules     interface{}
}

func (rg *RuleGroup) Evaluate(dataset map[string]interface{}) (bool, error) {
	rules := rg.Rules.([]interface{})

	switch rg.Condition.(string) {
	case AND:
		for _, r := range rules {
			isCorrectInput, err := rg.getChecker(r.(map[string]interface{})).Evaluate(dataset)
			if err != nil || !isCorrectInput {
				return isCorrectInput, err
			}
		}
		return true, nil

	case OR:
		for _, r := range rules {
			isCorrectInput, err := rg.getChecker(r.(map[string]interface{})).Evaluate(dataset)
			if err != nil || isCorrectInput {
				return isCorrectInput, err
			}
		}
		return false, nil

	default:
		return false, errors.Errorf("Invalid Condition %s", rg.Condition)
	}
}

func (rg *RuleGroup) getChecker(rule map[string]interface{}) Checker {
	if _, ok := rule["rules"]; ok {
		return &RuleGroup{Condition: rule["condition"], Rules: rule["rules"]}
	}

	r := &Rule{
		ID:       rule["id"].(string),
		Field:    rule["field"].(string),
		Type:     rule["type"].(string),
		Input:    rule["input"].(string),
		Operator: rule["operator"].(string),
		Value:    rule["value"],
	}

	if _, ok := rule["sanitize"]; ok {
		r.Sanitize = rule["sanitize"].(bool)
	}

	return r
}
