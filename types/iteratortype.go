package types

import (
	"io"

	"github.com/puppetlabs/go-evaluator/errors"
	"github.com/puppetlabs/go-evaluator/eval"
)

type (
	IteratorType struct {
		typ eval.PType
	}

	iteratorValue struct {
		iterator eval.Iterator
	}

	indexedIterator struct {
		elementType eval.PType
		pos         int
		indexed     eval.IndexedValue
	}

	mappingIterator struct {
		elementType eval.PType
		mapFunc     eval.Mapper
		base        eval.Iterator
	}

	predicateIterator struct {
		predicate eval.Predicate
		outcome   bool
		base      eval.Iterator
	}
)

var iteratorType_DEFAULT = &IteratorType{typ: DefaultAnyType()}

var Iterator_Type eval.ObjectType

func init() {
	Iterator_Type = newObjectType(`Pcore::IteratorType`,
		`Pcore::AnyType {
			attributes => {
				type => {
					type => Optional[Type],
					value => Any
				},
			}
		}`, func(ctx eval.EvalContext, args []eval.PValue) eval.PValue {
			return NewIteratorType2(args...)
		})
}

func DefaultIteratorType() *IteratorType {
	return iteratorType_DEFAULT
}

func NewIteratorType(elementType eval.PType) *IteratorType {
	if elementType == nil || elementType == anyType_DEFAULT {
		return DefaultIteratorType()
	}
	return &IteratorType{elementType}
}

func NewIteratorType2(args ...eval.PValue) *IteratorType {
	switch len(args) {
	case 0:
		return DefaultIteratorType()
	case 1:
		containedType, ok := args[0].(eval.PType)
		if !ok {
			panic(NewIllegalArgumentType2(`Iterator[]`, 0, `Type`, args[0]))
		}
		return NewIteratorType(containedType)
	default:
		panic(errors.NewIllegalArgumentCount(`Iterator[]`, `0 - 1`, len(args)))
	}
}

func (t *IteratorType) Accept(v eval.Visitor, g eval.Guard) {
	v(t)
	t.typ.Accept(v, g)
}

func (t *IteratorType) Default() eval.PType {
	return iteratorType_DEFAULT
}

func (t *IteratorType) Equals(o interface{}, g eval.Guard) bool {
	if ot, ok := o.(*IteratorType); ok {
		return t.typ.Equals(ot.typ, g)
	}
	return false
}

func (t *IteratorType) Generic() eval.PType {
	return NewIteratorType(eval.GenericType(t.typ))
}

func (t *IteratorType) Get(key string) (value eval.PValue, ok bool) {
	switch key {
	case `type`:
		return t.typ, true
	}
	return nil, false
}

func (t *IteratorType) IsAssignable(o eval.PType, g eval.Guard) bool {
	if it, ok := o.(*IteratorType); ok {
		return GuardedIsAssignable(t.typ, it.typ, g)
	}
	return false
}

func (t *IteratorType) IsInstance(o eval.PValue, g eval.Guard) bool {
	if it, ok := o.(eval.Iterator); ok {
		return GuardedIsInstance(t.typ, it.ElementType(), g)
	}
	return false
}

func (t *IteratorType) MetaType() eval.ObjectType {
	return Iterator_Type
}

func (t *IteratorType) Name() string {
	return `Iterator`
}

func (t *IteratorType) Parameters() []eval.PValue {
	if t.typ == DefaultAnyType() {
		return eval.EMPTY_VALUES
	}
	return []eval.PValue{t.typ}
}

func (t *IteratorType) String() string {
	return eval.ToString2(t, NONE)
}

func (t *IteratorType) ElementType() eval.PType {
	return t.typ
}

func (t *IteratorType) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	TypeToString(t, b, s, g)
}

func (t *IteratorType) Type() eval.PType {
	return &TypeType{t}
}

func WrapIterator(iter eval.Iterator) eval.IteratorValue {
	return &iteratorValue{iter}
}

func (it *iteratorValue) AsArray() eval.IndexedValue {
	return it.iterator.AsArray()
}

func (it *iteratorValue) Equals(o interface{}, g eval.Guard) bool {
	if ot, ok := o.(*iteratorValue); ok {
		return it.iterator.ElementType().Equals(ot.iterator.ElementType(), g)
	}
	return false
}

func (it *iteratorValue) Type() eval.PType {
	return NewIteratorType(it.iterator.ElementType())
}

func (it *iteratorValue) DynamicValue() eval.Iterator {
	return it.iterator
}

func (it *iteratorValue) String() string {
	return eval.ToString2(it, NONE)
}

func (it *iteratorValue) ToString(b io.Writer, s eval.FormatContext, g eval.RDetect) {
	if it.iterator.ElementType() != DefaultAnyType() {
		io.WriteString(b, `Iterator[`)
		eval.GenericType(it.iterator.ElementType()).ToString(b, s, g)
		io.WriteString(b, `]-Value`)
	} else {
		io.WriteString(b, `Iterator-Value`)
	}
}

func stopIteration() {
	if err := recover(); err != nil {
		if _, ok := err.(*errors.StopIteration); !ok {
			panic(err)
		}
	}
}

func find(iter eval.Iterator, predicate eval.Predicate, dflt eval.PValue, dfltProducer eval.Producer) (result eval.PValue) {
	defer stopIteration()

	result = eval.UNDEF
	ok := false
	for {
		result, ok = iter.Next()
		if !ok {
			if dfltProducer != nil {
				result = dfltProducer()
			} else {
				result = dflt
			}
			break
		}
		if predicate(result) {
			break
		}
	}
	return
}

func each(iter eval.Iterator, consumer eval.Consumer) {
	defer stopIteration()

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		consumer(v)
	}
}

func eachWithIndex(iter eval.Iterator, consumer eval.BiConsumer) {
	defer stopIteration()

	for idx := int64(0); ; idx++ {
		v, ok := iter.Next()
		if !ok {
			break
		}
		consumer(WrapInteger(idx), v)
	}
}

func all(iter eval.Iterator, predicate eval.Predicate) (result bool) {
	defer stopIteration()

	result = true
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if !predicate(v) {
			result = false
			break
		}
	}
	return
}

func any(iter eval.Iterator, predicate eval.Predicate) (result bool) {
	defer stopIteration()

	result = false
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if predicate(v) {
			result = true
			break
		}
	}
	return
}

func reduce2(iter eval.Iterator, value eval.PValue, redactor eval.BiMapper) (result eval.PValue) {
	defer stopIteration()

	result = value
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		result = redactor(result, v)
	}
	return
}

func reduce(iter eval.Iterator, redactor eval.BiMapper) eval.PValue {
	v, ok := iter.Next()
	if !ok {
		return _UNDEF
	}
	return reduce2(iter, v, redactor)
}

func asArray(iter eval.Iterator) (result eval.IndexedValue) {
	el := make([]eval.PValue, 0, 16)
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(*errors.StopIteration); ok {
				result = WrapArray(el)
			} else {
				panic(err)
			}
		}
	}()

	for {
		v, ok := iter.Next()
		if !ok {
			result = WrapArray(el)
			break
		}
		if it, ok := v.(eval.IteratorValue); ok {
			v = asArray(it.DynamicValue())
		}
		el = append(el, v)
	}
	return
}

func (ai *indexedIterator) All(predicate eval.Predicate) bool {
	return all(ai, predicate)
}

func (ai *indexedIterator) Any(predicate eval.Predicate) bool {
	return any(ai, predicate)
}

func (ai *indexedIterator) Each(consumer eval.Consumer) {
	each(ai, consumer)
}

func (ai *indexedIterator) EachWithIndex(consumer eval.BiConsumer) {
	eachWithIndex(ai, consumer)
}

func (ai *indexedIterator) ElementType() eval.PType {
	return ai.elementType
}

func (ai *indexedIterator) Find(predicate eval.Predicate) eval.PValue {
	return find(ai, predicate, _UNDEF, nil)
}

func (ai *indexedIterator) Find2(predicate eval.Predicate, dflt eval.PValue) eval.PValue {
	return find(ai, predicate, dflt, nil)
}

func (ai *indexedIterator) Find3(predicate eval.Predicate, dflt eval.Producer) eval.PValue {
	return find(ai, predicate, nil, dflt)
}

func (ai *indexedIterator) Next() (eval.PValue, bool) {
	pos := ai.pos + 1
	if pos < ai.indexed.Len() {
		ai.pos = pos
		return ai.indexed.At(pos), true
	}
	return _UNDEF, false
}

func (ai *indexedIterator) Map(elementType eval.PType, mapFunc eval.Mapper) eval.IteratorValue {
	return WrapIterator(&mappingIterator{elementType, mapFunc, ai})
}

func (ai *indexedIterator) Reduce(redactor eval.BiMapper) eval.PValue {
	return reduce(ai, redactor)
}

func (ai *indexedIterator) Reduce2(initialValue eval.PValue, redactor eval.BiMapper) eval.PValue {
	return reduce2(ai, initialValue, redactor)
}

func (ai *indexedIterator) Reject(predicate eval.Predicate) eval.IteratorValue {
	return WrapIterator(&predicateIterator{predicate, false, ai})
}

func (ai *indexedIterator) Select(predicate eval.Predicate) eval.IteratorValue {
	return WrapIterator(&predicateIterator{predicate, true, ai})
}

func (ai *indexedIterator) AsArray() eval.IndexedValue {
	return ai.indexed
}

func (ai *predicateIterator) All(predicate eval.Predicate) bool {
	return all(ai, predicate)
}

func (ai *predicateIterator) Any(predicate eval.Predicate) bool {
	return any(ai, predicate)
}

func (ai *predicateIterator) Next() (v eval.PValue, ok bool) {
	defer func() {
		if err := recover(); err != nil {
			if _, ok = err.(*errors.StopIteration); ok {
				ok = false
				v = _UNDEF
			} else {
				panic(err)
			}
		}
	}()

	for {
		v, ok = ai.base.Next()
		if !ok {
			v = _UNDEF
			break
		}
		if ai.predicate(v) == ai.outcome {
			break
		}
	}
	return
}

func (ai *predicateIterator) Each(consumer eval.Consumer) {
	each(ai, consumer)
}

func (ai *predicateIterator) EachWithIndex(consumer eval.BiConsumer) {
	eachWithIndex(ai, consumer)
}

func (ai *predicateIterator) ElementType() eval.PType {
	return ai.base.ElementType()
}

func (ai *predicateIterator) Find(predicate eval.Predicate) eval.PValue {
	return find(ai, predicate, _UNDEF, nil)
}

func (ai *predicateIterator) Find2(predicate eval.Predicate, dflt eval.PValue) eval.PValue {
	return find(ai, predicate, dflt, nil)
}

func (ai *predicateIterator) Find3(predicate eval.Predicate, dflt eval.Producer) eval.PValue {
	return find(ai, predicate, nil, dflt)
}

func (ai *predicateIterator) Map(elementType eval.PType, mapFunc eval.Mapper) eval.IteratorValue {
	return WrapIterator(&mappingIterator{elementType, mapFunc, ai})
}

func (ai *predicateIterator) Reduce(redactor eval.BiMapper) eval.PValue {
	return reduce(ai, redactor)
}

func (ai *predicateIterator) Reduce2(initialValue eval.PValue, redactor eval.BiMapper) eval.PValue {
	return reduce2(ai, initialValue, redactor)
}

func (ai *predicateIterator) Reject(predicate eval.Predicate) eval.IteratorValue {
	return WrapIterator(&predicateIterator{predicate, false, ai})
}

func (ai *predicateIterator) Select(predicate eval.Predicate) eval.IteratorValue {
	return WrapIterator(&predicateIterator{predicate, true, ai})
}

func (ai *predicateIterator) AsArray() eval.IndexedValue {
	return asArray(ai)
}

func (ai *mappingIterator) All(predicate eval.Predicate) bool {
	return all(ai, predicate)
}

func (ai *mappingIterator) Any(predicate eval.Predicate) bool {
	return any(ai, predicate)
}

func (ai *mappingIterator) Next() (v eval.PValue, ok bool) {
	v, ok = ai.base.Next()
	if !ok {
		v = _UNDEF
	} else {
		v = ai.mapFunc(v)
	}
	return
}

func (ai *mappingIterator) Each(consumer eval.Consumer) {
	each(ai, consumer)
}

func (ai *mappingIterator) EachWithIndex(consumer eval.BiConsumer) {
	eachWithIndex(ai, consumer)
}

func (ai *mappingIterator) ElementType() eval.PType {
	return ai.elementType
}

func (ai *mappingIterator) Find(predicate eval.Predicate) eval.PValue {
	return find(ai, predicate, _UNDEF, nil)
}

func (ai *mappingIterator) Find2(predicate eval.Predicate, dflt eval.PValue) eval.PValue {
	return find(ai, predicate, dflt, nil)
}

func (ai *mappingIterator) Find3(predicate eval.Predicate, dflt eval.Producer) eval.PValue {
	return find(ai, predicate, nil, dflt)
}

func (ai *mappingIterator) Map(elementType eval.PType, mapFunc eval.Mapper) eval.IteratorValue {
	return WrapIterator(&mappingIterator{elementType, mapFunc, ai})
}

func (ai *mappingIterator) Reduce(redactor eval.BiMapper) eval.PValue {
	return reduce(ai, redactor)
}

func (ai *mappingIterator) Reduce2(initialValue eval.PValue, redactor eval.BiMapper) eval.PValue {
	return reduce2(ai, initialValue, redactor)
}

func (ai *mappingIterator) Reject(predicate eval.Predicate) eval.IteratorValue {
	return WrapIterator(&predicateIterator{predicate, false, ai})
}

func (ai *mappingIterator) Select(predicate eval.Predicate) eval.IteratorValue {
	return WrapIterator(&predicateIterator{predicate, true, ai})
}

func (ai *mappingIterator) AsArray() eval.IndexedValue {
	return asArray(ai)
}
