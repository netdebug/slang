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
func constructHttpEndpoint(st storage.Storage, port int, opId uuid.UUID, gens core.Generics, props core.Properties) (*core.OperatorDef, error) {
	httpDef := &core.OperatorDef{
		Id:   "caff9fef-01fa-4ef8-bb11-aabbccddeeff",
		Meta: core.OperatorMetaDef{Name: "httpWrapper"},
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

	// The HTTP server is connected now, only the handler delegate is missing

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

	if inDef.Equals(elem.HTTP_REQUEST_DEF) {
		// If the operator can handle HTTP requests itself, just pass them
		httpDef.Connections["httpServer.handler)"] = []string{"(operator"}
	} else {
		// In this case we are not interested in anything but the body
		// It contains the JSON we need to unpack
		unpackerIns := &core.InstanceDef{
			Name:     "unpacker",
			Operator: elem.GetId("decode JSON").String(),
			Generics: core.Generics{
				"itemType": &inDef,
			},
		}
		httpDef.InstanceDefs = append(httpDef.InstanceDefs, unpackerIns)
		httpDef.Connections["httpServer.handler)body"] = []string{"(unpacker"}
		httpDef.Connections["unpacker)item"] = []string{"(operator"}
	}

	if outDef.Equals(elem.HTTP_RESPONSE_DEF) {
		// If the operator produces HTTP responses itself, just pass them
		httpDef.Connections["operator)"] = []string{"(httpServer.handler"}
	} else {
		// In this case we are not interested in anything but the body
		// It contains the JSON we need to pack
		packerIns := &core.InstanceDef{
			Name:     "packer",
			Operator: elem.GetId("encode JSON").String(),
			Generics: core.Generics{
				"itemType": &outDef,
			},
		}
		httpDef.InstanceDefs = append(httpDef.InstanceDefs, packerIns)
		httpDef.Connections["operator)"] = []string{"(packer"}
		// We connect unpacker output later

		// Now we still need status (200) and default headers ([])

		// Status code operator
		statusCodeIns := &core.InstanceDef{
			Name:     "statusCode",
			Operator: elem.GetId("value").String(),
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

		// Status code operator
		headersIns := &core.InstanceDef{
			Name:     "headers",
			Operator: elem.GetId("value").String(),
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

		httpDef.Connections["packer)"] = []string{"body(httpServer.handler", "(statusCode", "(headers"}
		httpDef.Connections["statusCode)"] = []string{"status(httpServer.handler"}
		httpDef.Connections["headers)"] = []string{"headers(httpServer.handler"}
	}

	return httpDef, nil
}
