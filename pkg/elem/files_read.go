package elem

import (
	"github.com/Bitspark/slang/pkg/core"
	"io/ioutil"
	"path/filepath"
	"github.com/Bitspark/slang/pkg/utils"
	"strings"
	"os/user"
)

var filesReadCfg = &builtinConfig{
	opDef: core.OperatorDef{
		ServiceDefs: map[string]*core.ServiceDef{
			core.MAIN_SERVICE: {
				In: core.TypeDef{
					Type: "string",
				},
				Out: core.TypeDef{
					Type: "map",
					Map: map[string]*core.TypeDef{
						"content": {
							Type: "binary",
						},
						"error": {
							Type: "string",
						},
					},
				},
			},
		},
	},
	opFunc: func(op *core.Operator) {
		in := op.Main().In()
		out := op.Main().Out()
		for !op.CheckStop() {
			file, marker := in.PullString()
			if marker != nil {
				out.Push(marker)
				continue
			}

			path := filepath.Clean(file)
			if strings.HasPrefix(path, "~") {
				usr, _ := user.Current()
				dir := usr.HomeDir
				path = filepath.Join(dir,path[1:])
			}
			content, err := ioutil.ReadFile(path)
			if err != nil {
				out.Map("content").Push(nil)
				out.Map("error").Push(err.Error())
				continue
			}

			out.Map("content").Push(utils.Binary(content))
			out.Map("error").Push(nil)
		}
	},
}