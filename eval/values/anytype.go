package values

import (
	. "io"

	. "github.com/puppetlabs/go-evaluator/eval/values/api"
)

type AnyType struct{}

func DefaultAnyType() *AnyType {
	return anyType_DEFAULT
}

func (t *AnyType) Equals(o interface{}, g Guard) bool {
	_, ok := o.(*AnyType)
	return ok
}

func (t *AnyType) IsAssignable(o PType, g Guard) bool {
	return true
}

func (t *AnyType) IsInstance(v PValue, g Guard) bool {
	return true
}

func (t *AnyType) Name() string {
	return `Any`
}

func (t *AnyType) String() string {
	return `Any`
}

func (t *AnyType) ToString(b Writer, s FormatContext, g RDetect) {
	TypeToString(t, b, s, g)
}

func (t *AnyType) Type() PType {
	return typeType_DEFAULT
}

var anyType_DEFAULT = &AnyType{}
