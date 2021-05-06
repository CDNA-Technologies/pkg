package querybuilder

import (
	"reflect"
	"strings"
	"sync"

	"github.com/enjoei/pkg/querybuilder/operator"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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
	var input, value interface{}
	var wg sync.WaitGroup
	var errg errgroup.Group

	opr, ok := operator.GetOperator(r.Operator)
	if !ok {
		return false, errors.Errorf("Invalid Operator %s", r.Operator)
	}

	wg.Add(2)

	errg.Go(func() error {
		var err error
		value, err = r.getValue()
		if err != nil {
			return err
		}
		wg.Done()
		return err
	})

	errg.Go(func() error {
		var err error
		input, err = r.getInputValue(dataset)
		if err != nil {
			return err
		}
		wg.Done()
		return err
	})

	if err := errg.Wait(); err != nil {
		return false, err
	}

	wg.Wait()
	return opr.Evaluate(input, value), nil
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
		if !ok || result == nil {
			return nil, errors.Errorf("Error in field: %s", field[i])
		}

		rresult := reflect.ValueOf(result)
		if rresult.Kind() == reflect.Map {
			rdataset = result.(map[string]interface{})
		} else if rresult.Kind() == reflect.Slice && i != (steps-1) {
			result = rresult.Index(0)
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
	rv := reflect.ValueOf(v)

	var errg errgroup.Group
	var wg sync.WaitGroup

	if rv.Kind() == reflect.Slice {
		sv := make([]interface{}, rv.Len())

		wg.Add(rv.Len())
		for _, vv := range v.([]interface{}) {
			go func(vv interface{}) {
				errg.Go(func() error {
					var err error
					_, err = r.castValue(vv)
					if err != nil {
						return err
					}
					return nil
				})
			}(vv)
			wg.Done()
		}

		if err := errg.Wait(); err != nil {
			return sv, err
		}

		wg.Wait()
		return sv, nil
	}

	return r.castValue(v)
}

// Available types in jQuery Query Builder are string, integer, double, date, time, datetime and boolean.
func (r *Rule) castValue(v interface{}) (interface{}, error) {

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
		return v, errors.Errorf("Invalid datatype: %s", r.Type)
	}
}
