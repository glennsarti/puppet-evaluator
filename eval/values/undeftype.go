package values

import (
	. "io"

	. "github.com/puppetlabs/go-evaluator/eval/values/api"
)

type (
	UndefType struct{}

	// UndefValue is an empty struct because both type and value are known
	UndefValue struct{}
)

var undefType_DEFAULT = &UndefType{}

func DefaultUndefType() *UndefType {
	return undefType_DEFAULT
}

func (t *UndefType) Equals(o interface{}, g Guard) bool {
	_, ok := o.(*UndefType)
	return ok
}

func (t *UndefType) IsAssignable(o PType, g Guard) bool {
	_, ok := o.(*UndefType)
	return ok
}

func (t *UndefType) IsInstance(o PValue, g Guard) bool {
	return o == _UNDEF
}

func (t *UndefType) Name() string {
	return `Undef`
}

func (t *UndefType) String() string {
	return `Undef`
}

func (t *UndefType) ToString(bld Writer, format FormatContext, g RDetect) {
	WriteString(bld, `Undef`)
}

func (t *UndefType) Type() PType {
	return &TypeType{t}
}

func WrapUndef() *UndefValue {
	return &UndefValue{}
}

func (uv *UndefValue) Equals(o interface{}, g Guard) bool {
	_, ok := o.(*UndefValue)
	return ok
}

func (uv *UndefValue) String() string {
	return `undef`
}

func (uv *UndefValue) ToKey() HashKey {
	return "\x01u"
}

func (uv *UndefValue) ToString(b Writer, s FormatContext, g RDetect) {
	WriteString(b, `undef`)
}

func (uv *UndefValue) Type() PType {
	return DefaultUndefType()
}
