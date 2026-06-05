package querybuilder

import (
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ruleNode interface {
	evaluate(dataset map[string]interface{}) evalResult
}

type evalResult struct {
	passed   bool
	failures []failure
	err      error
}

type failure struct {
	message  string
	priority int
}

type groupNode struct {
	condition string
	children  []ruleNode
}

func (g *groupNode) evaluate(dataset map[string]interface{}) evalResult {
	switch strings.ToUpper(g.condition) {
	case AND:
		return g.evaluateAND(dataset)
	case OR:
		return g.evaluateOR(dataset)
	default:
		return evalResult{err: errors.Errorf("unknown group condition: %q", g.condition)}
	}
}

func (g *groupNode) evaluateAND(dataset map[string]interface{}) evalResult {
	passed := true
	var failures []failure
	for _, child := range g.children {
		r := child.evaluate(dataset)
		if r.err != nil {
			return r
		}
		if !r.passed {
			passed = false
			failures = append(failures, r.failures...)
		}
	}
	return evalResult{passed: passed, failures: failures}
}

func (g *groupNode) evaluateOR(dataset map[string]interface{}) evalResult {
	var allFailures []failure
	for _, child := range g.children {
		r := child.evaluate(dataset)
		if r.err != nil {
			return r
		}
		if r.passed {
			return evalResult{passed: true}
		}
		allFailures = append(allFailures, r.failures...)
	}
	return evalResult{passed: false, failures: allFailures}
}

type leafNode struct {
	raw      map[string]interface{}
	message  string
	priority int
}

func (l *leafNode) evaluate(dataset map[string]interface{}) evalResult {
	wrapped := map[string]interface{}{
		"condition": AND,
		"rules":     []interface{}{l.raw},
	}
	qb := New(wrapped)
	passed, err := qb.Match(dataset)
	if err != nil {
		return evalResult{err: err}
	}
	if !passed {
		return evalResult{
			passed:   false,
			failures: []failure{{message: l.message, priority: l.priority}},
		}
	}
	return evalResult{passed: true}
}

// MatchWithFailureMessage evaluates the dataset against the ruleset in a single
// pass and returns the top-priority failure message when the match fails.
func (e *Evaluator) MatchWithFailureMessage(dataset map[string]interface{}) (bool, string, error) {
	result := buildTree(e.Ruleset).evaluate(dataset)
	if result.err != nil {
		return false, "", result.err
	}
	if result.passed {
		return true, "", nil
	}
	return false, topPriorityMessage(result.failures), nil
}

func buildTree(raw map[string]interface{}) ruleNode {
	if _, isGroup := raw["condition"]; isGroup {
		condition, _ := raw["condition"].(string)
		rulesRaw, _ := raw["rules"].([]interface{})
		children := make([]ruleNode, 0, len(rulesRaw))
		for _, r := range rulesRaw {
			if childMap, ok := r.(map[string]interface{}); ok {
				children = append(children, buildTree(childMap))
			}
		}
		return &groupNode{condition: condition, children: children}
	}

	message, _ := raw["validation_failed_message"].(string)
	priority := parsePriority(raw["validation_failed_message_priority"])
	return &leafNode{raw: raw, message: message, priority: priority}
}

func parsePriority(v interface{}) int {
	s, ok := v.(string)
	if !ok || s == "" {
		return 100
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 100
	}
	return n
}

func topPriorityMessage(failures []failure) string {
	if len(failures) == 0 {
		return ""
	}
	best := failures[0]
	for _, f := range failures[1:] {
		if effectivePriority(f.priority) < effectivePriority(best.priority) {
			best = f
		}
	}
	return best.message
}

func effectivePriority(p int) int {
	if p <= 0 {
		return math.MaxInt32
	}
	return p
}
