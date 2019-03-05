package functions

import (
	"github.com/lyraproj/pcore/eval"
	"github.com/lyraproj/pcore/types"
)

func init() {
	eval.NewGoFunction(`binary_file`,
		func(d eval.Dispatch) {
			d.Param(`String`)
			d.Function(func(c eval.Context, args []eval.Value) eval.Value {
				return types.BinaryFromFile(args[0].String())
			})
		})
}
