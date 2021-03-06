package functions

import (
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/puppet-evaluator/evaluator"
)

func init() {
	px.NewGoFunction(`flatten`,
		func(d px.Dispatch) {
			d.Param(`Iterable`)
			d.Function(func(c px.Context, args []px.Value) px.Value {
				switch arg := args[0].(type) {
				case px.List:
					return arg.Flatten()
				default:
					return evaluator.WrapIterable(arg.(px.Indexed)).AsArray().Flatten()
				}
			})
		},
	)
}
