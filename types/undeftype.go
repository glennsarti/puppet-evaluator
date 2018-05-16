package types

import (
	"io"

	"github.com/puppetlabs/go-evaluator/eval"
	"github.com/puppetlabs/go-issues/issue"
	"reflect"
)

type (
	UndefType struct{}

	// UndefValue is an empty struct because both type and value are known
	UndefValue struct{}
)

var undefType_DEFAULT = &UndefType{}

var Undef_Type eval.ObjectType

func init() {
	Undef_Type = newObjectType(`Pcore::UndefType`, `Pcore::AnyType{}`, func(ctx eval.Context, args []eval.PValue) eval.PValue {
		return DefaultUndefType()
	})
}

func DefaultUndefType() *UndefType {
	return undefType_DEFAULT
}

func (t *UndefType) Accept(v eval.Visitor, g eval.Guard) {
	v(t)
}

func (t *UndefType) Equals(o interface{}, g eval.Guard) bool {
	_, ok := o.(*UndefType)
	return ok
}

func (t *UndefType) IsAssignable(o eval.PType, g eval.Guard) bool {
	_, ok := o.(*UndefType)
	return ok
}

func (t *UndefType) IsInstance(c eval.Context, o eval.PValue, g eval.Guard) bool {
	return o == _UNDEF
}

func (t *UndefType) MetaType() eval.ObjectType {
	return Undef_Type
}

func (t *UndefType) Name() string {
	return `Undef`
}

func (t *UndefType) ReflectType() (reflect.Type, bool) {
	return reflect.Value{}.Type(), true
}

func (t *UndefType) String() string {
	return `Undef`
}

func (t *UndefType) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	TypeToString(t, b, s, g)
}

func (t *UndefType) Type() eval.PType {
	return &TypeType{t}
}

func WrapUndef() *UndefValue {
	return &UndefValue{}
}

func (uv *UndefValue) Equals(o interface{}, g eval.Guard) bool {
	_, ok := o.(*UndefValue)
	return ok
}

func (uv *UndefValue) Reflect(c eval.Context) reflect.Value {
	return reflect.Value{}
}

func (uv *UndefValue) ReflectTo(c eval.Context, value reflect.Value) {
	if !value.CanSet() {
		panic(eval.Error(c, eval.EVAL_ATTEMPT_TO_SET_UNSETTABLE, issue.H{`kind`: value.Kind().String()}))
	}
	value.Set(reflect.Zero(value.Type()))
}

func (uv *UndefValue) String() string {
	return `undef`
}

func (uv *UndefValue) ToKey() eval.HashKey {
	return eval.HashKey([]byte{1, HK_UNDEF})
}

func (uv *UndefValue) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	io.WriteString(b, `undef`)
}

func (uv *UndefValue) Type() eval.PType {
	return DefaultUndefType()
}
