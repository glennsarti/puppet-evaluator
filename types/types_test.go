package types_test

import (
	"fmt"
	"os"
	"regexp"

	"github.com/puppetlabs/go-evaluator/eval"
	"github.com/puppetlabs/go-evaluator/types"

	// Ensure that pcore is initialized
	_ "github.com/puppetlabs/go-evaluator/pcore"
)

func ExampleUniqueValues() {
	x := types.WrapString(`hello`)
	y := types.WrapInteger(32)
	types.UniqueValues([]eval.Value{x, y})

	z := types.WrapString(`hello`)
	svec := []*types.StringValue{x, z}
	fmt.Println(types.UniqueValues([]eval.Value{svec[0], svec[1]}))
	// Output: [hello]
}

func ExampleNewCallableType2() {
	cc := types.NewCallableType2(types.ZERO, types.WrapDefault())
	fmt.Println(cc)
	// Output: Callable[0, default]
}

func ExampleNewTupleType() {
	tuple := types.NewTupleType([]eval.Type{types.DefaultStringType(), types.DefaultIntegerType()}, nil)
	fmt.Println(tuple)
	// Output: Tuple[String, Integer]
}

func ExampleWrapHash() {
	a := eval.Wrap(nil, map[string]interface{}{
		`foo`: 23,
		`fee`: `hello`,
		`fum`: map[string]interface{}{
			`x`: `1`,
			`y`: []int{1, 2, 3},
			`z`: regexp.MustCompile(`^[a-z]+$`)}})

	e := types.WrapHash([]*types.HashEntry{
		types.WrapHashEntry2(`foo`, types.WrapInteger(23)),
		types.WrapHashEntry2(`fee`, types.WrapString(`hello`)),
		types.WrapHashEntry2(`fum`, types.WrapHash([]*types.HashEntry{
			types.WrapHashEntry2(`x`, types.WrapString(`1`)),
			types.WrapHashEntry2(`y`, types.WrapValues([]eval.Value{
				types.WrapInteger(1), types.WrapInteger(2), types.WrapInteger(3)})),
			types.WrapHashEntry2(`z`, types.WrapRegexp(`^[a-z]+$`))}))})

	fmt.Println(eval.Equals(e, a))
	// Output: true
}

func ExampleNew() {
	eval.Puppet.Do(func(c eval.Context) {
		t := c.ParseType2(`Struct[{'name' => String, 'type' => Type, 'value' => Variant[Deferred, Data], Optional['captures_rest'] => Boolean}]`)
		t.ToString(os.Stdout, eval.PRETTY, nil)
		fmt.Println()
	})
	// Output: hello
}