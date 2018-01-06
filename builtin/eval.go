package builtin

import (
	"errors"
	"slang/core"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

type EvaluableExpression struct {
	govaluate.EvaluableExpression
}

func getFlattenedObj(obj interface{}) map[string]interface{} {
	flatMap := make(map[string]interface{})

	if a, ok := obj.([]interface{}); ok {
		for k, val := range a {
			key := strconv.Itoa(k)
			var sub interface{}
			var ok bool

			if sub, ok = val.(map[string]interface{}); !ok {
				if sub, ok = val.([]interface{}); !ok {
					flatMap[key] = val
					continue
				}
			}

			for sKey, sVal := range getFlattenedObj(sub) {
				flatMap[key+"__"+sKey] = sVal
			}
		}

	} else if m, ok := obj.(map[string]interface{}); ok {
		for key, val := range m {
			var sub interface{}
			var ok bool

			if sub, ok = val.(map[string]interface{}); !ok {
				if sub, ok = val.([]interface{}); !ok {
					flatMap[key] = val
					continue
				}
			}

			for sKey, sVal := range getFlattenedObj(sub) {
				flatMap[key+"__"+sKey] = sVal
			}
		}

	} else {
		panic("obj must be list or map")
	}

	return flatMap
}

func NewFlatMapParameters(m map[string]interface{}) govaluate.MapParameters {
	flatMap := getFlattenedObj(m)
	return govaluate.MapParameters(flatMap)
}

func NewEvaluableExpression(expression string) (*EvaluableExpression, error) {
	goEvalExpr, err := govaluate.NewEvaluableExpression(strings.Replace(expression, ".", "__", -1))
	if err == nil {
		return &EvaluableExpression{*goEvalExpr}, nil
	}
	return nil, err
}

type functionStore struct {
	expr     string
	evalExpr *govaluate.EvaluableExpression
}

func createOpEval(def core.InstanceDef, par *core.Operator) (*core.Operator, error) {
	if def.Properties == nil {
		return nil, errors.New("no properties given")
	}

	exprStr, ok := def.Properties["expression"]

	if !ok {
		return nil, errors.New("no expression given")
	}

	expr, ok := exprStr.(string)

	if !ok {
		return nil, errors.New("expression must be string")
	}

	evalExpr, err := govaluate.NewEvaluableExpression(expr)

	if err != nil {
		return nil, err
	}

	inDef := core.PortDef{
		Type: "map",
		Map:  make(map[string]core.PortDef),
	}

	vars := evalExpr.Vars()

	for _, v := range vars {
		inDef.Map[v] = core.PortDef{Type: "any"}
	}

	outDef := core.PortDef{
		Type: "any",
	}

	o, err := core.NewOperator(def.Name, func(in, out *core.Port, store interface{}) {
		expr := store.(functionStore).evalExpr
		for true {
			i := in.Pull()

			if isMarker(i) {
				out.Push(i)
				continue
			}

			if m, ok := i.(map[string]interface{}); ok {
				rlt, _ := expr.Eval(NewFlatMapParameters(m))
				out.Push(rlt)
			} else {
				panic("invalid item")
			}
		}
	}, inDef, outDef, par)
	o.SetStore(functionStore{expr, evalExpr})

	return o, nil
}
