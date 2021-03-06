package daemon

import (
	"github.com/Bitspark/slang/pkg/api"
	"github.com/Bitspark/slang/pkg/core"
	"github.com/Bitspark/slang/pkg/elem"
	"github.com/Bitspark/slang/pkg/storage"
	"github.com/google/uuid"
)

// Constructs an executable operator
// TODO: Make safer (maybe require an API key?)
func constructHttpStreamEndpoint(st storage.Storage, port int, opId uuid.UUID, gens core.Generics, props core.Properties) (*core.OperatorDef, error) {
	httpDef := &core.OperatorDef{
		Id: "caff9fef-01fa-4ef8-bb11-aabbccddeeff",
		ServiceDefs: map[string]*core.ServiceDef{
			core.MAIN_SERVICE: {
				In: core.TypeDef{
					Type: "trigger",
				},
				Out: core.TypeDef{
					Type: "trigger",
				},
			},
		},
		Connections: make(map[string][]string),
	}

	op, err := api.Build(opId, gens, props, st)
	if err != nil {
		return nil, err
	}

	// Const port instance
	portIns := &core.InstanceDef{
		Name:     "port",
		Operator: elem.GetId("value").String(),
		Generics: core.Generics{
			"valueType": {
				Type: "number",
			},
		},
		Properties: core.Properties{
			"value": float64(port),
		},
	}
	httpDef.InstanceDefs = append(httpDef.InstanceDefs, portIns)
	httpDef.Connections["("] = []string{"(port"}

	// HTTP operator instance
	httpIns := &core.InstanceDef{
		Name:     "httpServer",
		Operator: elem.GetId("HTTP server").String(),
	}
	httpDef.InstanceDefs = append(httpDef.InstanceDefs, httpIns)
	httpDef.Connections["port)"] = []string{"(httpServer"}
	httpDef.Connections["httpServer)"] = []string{")"}

	// The HTTP server is connected now, the handler delegate however is missing

	// We need a switch to determine if a new value is pushed or a value request is made
	switchIns := &core.InstanceDef{
		Name:     "switch",
		Operator: "slang.control.Switch",
		Generics: map[string]*core.TypeDef{
			"inType": {
				Type: "binary",
			},
			"outType": {
				Type: "binary",
			},
			"selectType": {
				Type: "string",
			},
		},
		Properties: map[string]interface{}{
			"cases": []interface{}{"POST"},
		},
	}
	httpDef.InstanceDefs = append(httpDef.InstanceDefs, switchIns)
	httpDef.Connections["httpServer.handler)method"] = []string{"select(switch"}
	httpDef.Connections["httpServer.handler)body"] = []string{"item(switch"}

	// Status code operator
	statusCodeIns := &core.InstanceDef{
		Name:     "statusCode",
		Operator: "slang.data.Value",
		Generics: core.Generics{
			"valueType": {
				Type: "number",
			},
		},
		Properties: core.Properties{
			"value": 200,
		},
	}
	httpDef.InstanceDefs = append(httpDef.InstanceDefs, statusCodeIns)
	// We connect it later

	// Header operator
	headersIns := &core.InstanceDef{
		Name:     "headers",
		Operator: "slang.data.Value",
		Generics: core.Generics{
			"valueType": {
				Type: "stream",
				Stream: &core.TypeDef{
					Type: "map",
					Map: map[string]*core.TypeDef{
						"key": {
							Type: "string",
						},
						"value": {
							Type: "string",
						},
					},
				},
			},
		},
		Properties: core.Properties{
			"value": []interface{}{
				map[string]interface{}{"key": "Access-Control-Allow-Origin", "value": "*"},
				map[string]interface{}{"key": "Content-Type", "value": "application/json"},
			},
		},
	}
	httpDef.InstanceDefs = append(httpDef.InstanceDefs, headersIns)
	// We connect it later

	httpDef.Connections["switch)"] = []string{"body(httpServer.handler", "(statusCode", "(headers"}
	httpDef.Connections["statusCode)"] = []string{"status(httpServer.handler"}
	httpDef.Connections["headers)"] = []string{"headers(httpServer.handler"}

	// Now the switch delegates are still missing

	// This is the actual operator we want to execute
	operatorIns := &core.InstanceDef{
		Name:       "operator",
		Operator:   opId.String(),
		Generics:   gens,
		Properties: props,
	}
	httpDef.InstanceDefs = append(httpDef.InstanceDefs, operatorIns)

	// Get operator interface
	inDef := op.Main().In().Define()
	outDef := op.Main().Out().Define()

	// This is the store holding all data
	storeIns := &core.InstanceDef{
		Name:     "store",
		Operator: "slang.meta.Store",
		Generics: map[string]*core.TypeDef{
			"examineType": &outDef,
		},
		Properties: props,
	}
	httpDef.InstanceDefs = append(httpDef.InstanceDefs, storeIns)

	// POST case

	unpackerIns := &core.InstanceDef{
		Name:     "unpacker",
		Operator: "slang.encoding.JSONRead",
		Generics: core.Generics{
			"itemType": &inDef,
		},
	}
	httpDef.InstanceDefs = append(httpDef.InstanceDefs, unpackerIns)
	httpDef.Connections["switch.POST)"] = []string{"(unpacker", "(switch.POST"}
	httpDef.Connections["unpacker)item"] = []string{"(operator"}
	httpDef.Connections["operator)"] = []string{"(store"}

	// GET (default) case

	packerIns := &core.InstanceDef{
		Name:     "packer",
		Operator: "slang.encoding.JSONWrite",
		Generics: core.Generics{
			"itemType": &core.TypeDef{
				Type:   "stream",
				Stream: &outDef,
			},
		},
	}
	httpDef.InstanceDefs = append(httpDef.InstanceDefs, packerIns)
	httpDef.Connections["switch.default)"] = []string{"(query@store"}
	httpDef.Connections["query@store)"] = []string{"(packer"}
	httpDef.Connections["packer)"] = []string{"(switch.default"}

	return httpDef, nil
}
