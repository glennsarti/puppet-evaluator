package values

import (
	. "io"

	. "github.com/puppetlabs/go-evaluator/eval/values/api"
)

type NumericType struct{}

func DefaultNumericType() *NumericType {
	return numericType_DEFAULT
}
func (t *NumericType) Equals(o interface{}, g Guard) bool {
	_, ok := o.(*NumericType)
	return ok
}

func (t *NumericType) IsAssignable(o PType, g Guard) bool {
	switch o.(type) {
	case *IntegerType, *FloatType:
		return true
	default:
		return false
	}
}

func (t *NumericType) IsInstance(o PValue, g Guard) bool {
	switch o.Type().(type) {
	case *FloatType, *IntegerType:
		return true
	default:
		return false
	}
}

func (t *NumericType) Name() string {
	return `Numeric`
}

func (t *NumericType) String() string {
	return ToString2(t, NONE)
}

func (t *NumericType) ToString(bld Writer, format FormatContext, g RDetect) {
	WriteString(bld, `Numeric`)
}

func (t *NumericType) Type() PType {
	return &TypeType{t}
}

var numericType_DEFAULT = &NumericType{}
