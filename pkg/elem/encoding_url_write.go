package elem

import (
	"github.com/Bitspark/slang/pkg/core"
	"net/url"
	"fmt"
)

var encodingURLWriteCfg = &builtinConfig{
	opDef: core.OperatorDef{
		Id: "702a2036-a1cc-4783-8b83-b18494c5e9f1",
		Meta: core.OperatorMetaDef{
			Name: "encode URL",
			ShortDescription: "encodes a Slang map into the corresponding URL-encoded string",
			Icon: "brackets",
			Tags: []string{"http", "encoding"},
			DocURL: "https://bitspark.de/slang/docs/operator/encode-url",
		},
		ServiceDefs: map[string]*core.ServiceDef{
			core.MAIN_SERVICE: {
				In: core.TypeDef{
					Type: "map",
					Map: map[string]*core.TypeDef{
						"{params}": {
							Type: "primitive",
						},
					},
				},
				Out: core.TypeDef{
					Type: "string",
				},
			},
		},
		DelegateDefs: map[string]*core.DelegateDef{},
		PropertyDefs: map[string]*core.TypeDef{
			"params": {
				Type: "stream",
				Stream: &core.TypeDef{
					Type: "string",
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
			vals := url.Values{}
			im := i.(map[string]interface{})
			for key, value := range im {
				vals.Set(key, fmt.Sprintf("%v", value))
			}
			out.Push(vals.Encode())
		}
	},
}
