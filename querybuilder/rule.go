package querybuilder

import (
	"reflect"
	"strings"
	"sync"

	"github.com/enjoei/pkg/querybuilder/operator"
)

type Rule struct {
	ID       string
	Field    string
	Type     string
	Input    string
	Operator string
	Sanitize bool
	Value    interface{}
}

// Evaluate function checks whether the dataset matches with rule
func (r *Rule) Evaluate(dataset map[string]interface{}) (bool, error) {
	var err error = nil
	var wg sync.WaitGroup
	var input, value interface{}

	opr, ok := operator.GetOperator(r.Operator)
	if !ok {
		return false, err
	}

	wg.Add(2)
	go func() {
		input, err = r.getInputValue(dataset)
		wg.Done()
	}()

	go func() {
		value, err = r.getValue()
		wg.Done()
	}()

	wg.Wait()

	if err != nil {
		return false, err
	}

	return opr.Evaluate(input, value), err
}

func (r *Rule) getValue() (interface{}, error) {
	return r.parseValue(r.Value)
}

// getInputValue fetch in the dataset the field value and convert to the type of the rule
func (r *Rule) getInputValue(dataset map[string]interface{}) (interface{}, error) {
	var rdataset = make(map[string]interface{})
	var result interface{}
	var ok bool

	for k, v := range dataset {
		rdataset[k] = v
	}

	field := strings.Split(r.Field, ".")
	steps := len(field)

	for i := 0; i < steps; i++ {
		result, ok = rdataset[field[i]]
		if !ok {
			return nil, nil
		}

		rresult := reflect.ValueOf(result)
		if rresult.Kind() == reflect.Map {
			rdataset = result.(map[string]interface{})
		} else if rresult.Kind() == reflect.Slice && i != (steps-1) {
			result = rresult.Index(0)
		}

		if result == nil {
			return nil, nil
		}
	}

	iv, err := r.parseValue(result)

	if r.Sanitize && r.Type == "string" {
		v := iv.(string)
		return sanitize(&v), err
	}

	return iv, err
}

func (r *Rule) parseValue(v interface{}) (interface{}, error) {
	var err error = nil

	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Slice {
		sv := make([]interface{}, rv.Len())

		for i, vv := range v.([]interface{}) {
			sv[i], err = r.castValue(vv)
		}
		return sv, err
	}

	return r.castValue(v)
}

// Available types in jQuery Query Builder are string, integer, double, date, time, datetime and boolean.
func (r *Rule) castValue(v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}

	switch r.Type {
	case "string":
		return toString(v)
	case "integer":
		return toInteger(v)
	case "double":
		return toDouble(v)
	case "date":
		return toDate(v)
	case "time":
		return toTime(v)
	case "datetime":
		return toDateTime(v)
	case "boolean":
		return toBoolean(v)
	default:
		return v, nil
	}
}
