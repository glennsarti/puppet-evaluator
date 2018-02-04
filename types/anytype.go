package types

import (
	"io"

	"github.com/puppetlabs/go-evaluator/eval"
)

type AnyType struct{}

var Any_Type eval.ObjectType

func init() {
	Any_Type = newObjectType(`Pcore::AnyType`, `{}`, func(ctx eval.EvalContext, args []eval.PValue) eval.PValue {
		return DefaultAnyType()
	})
}

func DefaultAnyType() *AnyType {
	return anyType_DEFAULT
}

func (t *AnyType) Accept(v eval.Visitor, g eval.Guard) {
	v(t)
}

func (t *AnyType) Equals(o interface{}, g eval.Guard) bool {
	_, ok := o.(*AnyType)
	return ok
}

func (t *AnyType) IsAssignable(o eval.PType, g eval.Guard) bool {
	return true
}

func (t *AnyType) IsInstance(v eval.PValue, g eval.Guard) bool {
	return true
}

func (t *AnyType) MetaType() eval.ObjectType {
	return Any_Type
}

func (t *AnyType) Name() string {
	return `Any`
}

func (t *AnyType) String() string {
	return `Any`
}

func (t *AnyType) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	TypeToString(t, b, s, g)
}

func (t *AnyType) Type() eval.PType {
	return &TypeType{t}
}

var anyType_DEFAULT = &AnyType{}
