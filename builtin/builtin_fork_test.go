package builtin

import (
	"slang/op"
	"slang/tests"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperatorCreator_Fork_IsRegistered(t *testing.T) {
	a := assert.New(t)

	ocFork := getCreatorFunc("fork")
	a.NotNil(ocFork)
}

func TestBuiltin_OperatorFork__InPorts(t *testing.T) {
	a := assert.New(t)

	o, err := getCreatorFunc("fork")(op.InstanceDef{Operator: "fork"}, nil)
	a.NoError(err)

	a.NotNil(o.In().Stream().Map("i"))
	a.NotNil(o.In().Stream().Map("select"))
}

func TestBuiltin_OperatorFork__OutPorts(t *testing.T) {
	a := assert.New(t)

	o, err := getCreatorFunc("fork")(op.InstanceDef{Operator: "fork"}, nil)
	a.NoError(err)

	a.NotNil(o.Out().Map("true").Stream())
	a.NotNil(o.Out().Map("false").Stream())
}

func TestBuiltin_OperatorFork__Correct(t *testing.T) {
	a := assert.New(t)

	o, err := getCreatorFunc("fork")(op.InstanceDef{Operator: "fork"}, nil)
	a.NoError(err)

	o.Out().Map("true").Stream().Bufferize()
	o.Out().Map("false").Stream().Bufferize()
	o.Start()

	o.In().Push([]interface{}{
		map[string]interface{}{
			"i":      "hallo",
			"select": true,
		},
		map[string]interface{}{
			"i":      "welt",
			"select": false,
		},
		map[string]interface{}{
			"i":      100,
			"select": true,
		},
		map[string]interface{}{
			"i":      101,
			"select": false,
		},
	})

	tests.AssertPortItems(t, []interface{}{[]interface{}{"hallo", 100}}, o.Out().Map("true"))
	tests.AssertPortItems(t, []interface{}{[]interface{}{"welt", 101}}, o.Out().Map("false"))
}
