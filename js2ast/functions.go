package js2ast

import (
	"github.com/puppetlabs/go-evaluator/eval"
	"github.com/puppetlabs/go-evaluator/types"
)

func eachIterator(c eval.Context, arg eval.IterableValue, block eval.Lambda) eval.Value {
	arg.Iterator().Each(func(v eval.Value) { block.Call(c, nil, v) })
	return eval.UNDEF
}

func InitJavaScript(c eval.Context) {
	if c.Language() == eval.LangJavaScript {
		return
	}
	c.SetLanguage(eval.LangJavaScript)

	initConsole(c)

	eval.NewGoFunction(`js::forIn`,
		func(d eval.Dispatch) {
			d.Param(`Variant[Array,String]`)
			d.Block(`Callable[1,1]`)
			d.Function2(func(c eval.Context, args []eval.Value, block eval.Lambda) eval.Value {
				return eachIterator(c, args[0].(eval.List), block)
			})
		},

		func(d eval.Dispatch) {
			d.Param(`Hash`)
			d.Block(`Callable[1,1]`)
			d.Function2(func(c eval.Context, args []eval.Value, block eval.Lambda) eval.Value {
				return eachIterator(c, args[0].(*types.HashValue).Keys(), block)
			})
		})

	// NOTE: This is not a proper implementation of JavaScript instanceof since the
	// evaluator uses types instead of prototypes.
	//
	// A 100% compliant implementation must yield the somewhat peculiar result:
	//
	// ('hello' instanceof String)             // false
	// (new String('hello') instanceof String) // true
  //
	eval.NewGoFunction(`js::isInstance`,
		func(d eval.Dispatch) {
			d.Param(`Type`)
			d.Param(`Any`)
			d.Function(func(c eval.Context, args []eval.Value) eval.Value {
				return types.WrapBoolean(eval.IsInstance(args[0].(eval.Type), args[1]))
			})
		})

	c.ResolveResolvables()
	c.Scope().Set(`console`, NewConsole(`default`))
}
