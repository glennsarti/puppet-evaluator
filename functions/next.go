package functions

import (
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/puppet-evaluator/errors"
)

func init() {
	px.NewGoFunction(`next`,
		func(d px.Dispatch) {
			d.OptionalParam(`Any`)
			d.Function(func(c px.Context, args []px.Value) px.Value {
				arg := px.Undef
				if len(args) > 0 {
					arg = args[0]
				}
				panic(errors.NewNextIteration(c.StackTop(), arg))
			})
		},
	)
}
