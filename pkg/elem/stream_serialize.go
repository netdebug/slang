package elem

import (
	"github.com/Bitspark/slang/pkg/core"
	"strings"
	"strconv"
)

var streamSerializeCfg = &builtinConfig{
	opDef: core.OperatorDef{
		Id: "13257172-b05d-497c-be23-da7c86577c1e",
		Meta: core.OperatorMetaDef{
			Name: "serialize",
			ShortDescription: "takes a map of items and serializes them into a stream",
			Icon: "ellipsis-h",
			Tags: []string{"stream", "convert"},
			DocURL: "https://bitspark.de/slang/docs/operator/serialize",
		},
		ServiceDefs: map[string]*core.ServiceDef{
			core.MAIN_SERVICE: {
				In: core.TypeDef{
					Type: "map",
					Map: map[string]*core.TypeDef{
						"el_{indexes}": {
							Type:    "generic",
							Generic: "itemType",
						},
					},
				},
				Out: core.TypeDef{
					Type: "stream",
					Stream: &core.TypeDef{
						Type:    "generic",
						Generic: "itemType",
					},
				},
			},
		},
		DelegateDefs: map[string]*core.DelegateDef{},
		PropertyDefs: map[string]*core.TypeDef{
			"indexes": {
				Type: "stream",
				Stream: &core.TypeDef{
					Type: "number",
				},
			},
		},
	},
	opFunc: func(op *core.Operator) {
		in := op.Main().In()
		out := op.Main().Out()
		for !op.CheckStop() {
			i := in.Pull()
			if core.IsMarker(i) {
				out.Push(i)
				continue
			}

			im := i.(map[string]interface{})

			maxIndex := -1
			for k := range im {
				if !strings.HasPrefix(k, "el_") {
					panic("malformed entry")
				}
				index, _ := strconv.Atoi(k[3:])
				if index > maxIndex {
					maxIndex = index
				}
			}

			stream := make([]interface{}, maxIndex + 1)
			for k, v := range im {
				index, _ := strconv.Atoi(k[3:])
				stream[index] = v
			}

			out.Push(stream)
		}
	},
}
