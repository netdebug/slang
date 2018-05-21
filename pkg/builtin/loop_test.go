package builtin

import (
	"github.com/Bitspark/slang/pkg/core"
	"github.com/Bitspark/slang/tests/assertions"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestOperatorCreator_Loop_IsRegistered(t *testing.T) {
	a := assertions.New(t)

	ocLoop := getBuiltinCfg("slang.loop")
	a.NotNil(ocLoop)
}

func TestBuiltin_Loop__Simple(t *testing.T) {
	a := assertions.New(t)
	lo, err := buildOperator(
		core.InstanceDef{
			Name:     "loop",
			Operator: "slang.loop",
			Generics: map[string]*core.TypeDef{
				"stateType": {
					Type: "number",
				},
			},
		},
	)
	a.NoError(err)
	a.NotNil(lo)

	// Condition operator
	co, _ := core.NewOperator(
		"cond",
		func(op *core.Operator) {
			in := op.Main().In()
			out := op.Main().Out()
			for {
				i := in.Pull()
				f, ok := i.(float64)
				if !ok {
					out.Push(i)
				} else {
					out.Push(f < 10.0)
				}
			}
		},
		nil,
		nil,
		nil,
		core.OperatorDef{
			ServiceDefs: map[string]*core.ServiceDef{"main": {In: core.TypeDef{Type: "number"}, Out: core.TypeDef{Type: "boolean"}}},
		},
	)

	// Double function operator
	fo, _ := core.NewOperator(
		"double",
		func(op *core.Operator) {
			in := op.Main().In()
			out := op.Main().Out()
			for {
				i := in.Pull()
				f, ok := i.(float64)
				if !ok {
					out.Push(i)
				} else {
					out.Push(f * 2.0)
				}
			}
		},
		nil,
		nil,
		nil,
		core.OperatorDef{
			ServiceDefs: map[string]*core.ServiceDef{"main": {In: core.TypeDef{Type: "number"}, Out: core.TypeDef{Type: "number"}}},
		},
	)

	// Connect
	a.NoError(lo.Delegate("iteration").Out().Stream().Connect(fo.Main().In()))
	a.NoError(lo.Delegate("iteration").Out().Stream().Connect(co.Main().In()))
	a.NoError(fo.Main().Out().Connect(lo.Delegate("iteration").In().Stream().Map("state")))
	a.NoError(co.Main().Out().Connect(lo.Delegate("iteration").In().Stream().Map("continue")))

	lo.Main().Out().Bufferize()

	lo.Main().In().Push(1.0)
	lo.Main().In().Push(10.0)

	lo.Start()
	fo.Start()
	co.Start()

	a.PortPushesAll([]interface{}{16.0, 10.0}, lo.Main().Out())
}

func
TestBuiltin_Loop__Fibo(t *testing.T) {
	a := assertions.New(t)
	stateType := core.TypeDef{
		Type: "map",
		Map: map[string]*core.TypeDef{
			"i":      {Type: "number"},
			"fib":    {Type: "number"},
			"oldFib": {Type: "number"},
		},
	}
	lo, err := buildOperator(
		core.InstanceDef{
			Operator: "slang.loop",
			Generics: map[string]*core.TypeDef{
				"stateType": &stateType,
			},
		},
	)
	require.NoError(t, err)
	a.NotNil(lo)
	require.Equal(t, core.TYPE_MAP, lo.Main().In().Type())
	require.Equal(t, core.TYPE_NUMBER, lo.Main().In().Map("i").Type())

	// Condition operator
	co, _ := core.NewOperator(
		"cond",
		func(op *core.Operator) {
			in := op.Main().In()
			out := op.Main().Out()
			for {
				i := in.Pull()
				fm, ok := i.(map[string]interface{})
				if !ok {
					out.Push(i)
				} else {
					i := fm["i"].(float64)
					out.Push(i > 0.0)
				}
			}
		},
		nil,
		nil,
		nil,
		core.OperatorDef{
			ServiceDefs: map[string]*core.ServiceDef{"main": {In: stateType, Out: core.TypeDef{Type: "boolean"}}},
		},
	)

	// Fibonacci function operator
	fo, _ := core.NewOperator(
		"fib",
		func(op *core.Operator) {
			in := op.Main().In()
			out := op.Main().Out()
			for {
				i := in.Pull()
				fm, ok := i.(map[string]interface{})
				if !ok {
					out.Push(i)
				} else {
					i := fm["i"].(float64) - 1
					oldFib := fm["fib"].(float64)
					fib := fm["oldFib"].(float64) + oldFib
					out.Push(map[string]interface{}{"i": i, "fib": fib, "oldFib": oldFib})
				}
			}
		},
		nil,
		nil,
		nil,
		core.OperatorDef{
			ServiceDefs: map[string]*core.ServiceDef{"main": {In: stateType, Out: stateType}},
		},
	)

	// Connect
	a.NoError(lo.Delegate("iteration").Out().Stream().Connect(fo.Main().In()))
	a.NoError(lo.Delegate("iteration").Out().Stream().Connect(co.Main().In()))
	a.NoError(fo.Main().Out().Connect(lo.Delegate("iteration").In().Stream().Map("state")))
	a.NoError(co.Main().Out().Connect(lo.Delegate("iteration").In().Stream().Map("continue")))

	lo.Main().Out().Bufferize()

	lo.Main().In().Push(map[string]interface{}{"i": 10.0, "fib": 1.0, "oldFib": 0.0})
	lo.Main().In().Push(map[string]interface{}{"i": 20.0, "fib": 1.0, "oldFib": 0.0})

	lo.Start()
	fo.Start()
	co.Start()

	a.PortPushesAll([]interface{}{
		map[string]interface{}{"i": 0.0, "fib": 89.0, "oldFib": 55.0},
		map[string]interface{}{"i": 0.0, "fib": 10946.0, "oldFib": 6765.0},
	}, lo.Main().Out())
}

func
TestBuiltin_Loop__MarkersPushedCorrectly(t *testing.T) {
	a := assertions.New(t)
	lo, err := buildOperator(
		core.InstanceDef{
			Operator: "slang.loop",
			Generics: map[string]*core.TypeDef{
				"stateType": {
					Type: "number",
				},
			},
		},
	)
	a.NoError(err)
	a.NotNil(lo)

	lo.Main().Out().Bufferize()
	lo.Delegate("iteration").Out().Bufferize()

	lo.Start()

	pInit := lo.Main().In()
	pIteration := lo.Delegate("iteration").In()
	pState := lo.Delegate("iteration").Out().Stream()
	pEnd := lo.Main().Out()

	bos := core.BOS{}
	pInit.Push(bos)
	a.Nil(pEnd.Poll())
	a.Equal(bos, pState.Pull())

	pIteration.Push(bos)

	a.Equal(bos, pEnd.Pull())
	a.Nil(pState.Poll())

	eos := core.BOS{}
	pInit.Push(eos)
	a.Nil(pEnd.Poll())
	a.Equal(eos, pState.Pull())

	pIteration.Push(eos)

	a.Equal(eos, pEnd.Pull())
	a.Nil(pState.Poll())
}