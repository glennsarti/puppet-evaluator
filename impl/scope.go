package impl

import (
	"strings"

	"github.com/puppetlabs/go-evaluator/eval"
	"github.com/puppetlabs/go-evaluator/types"
)

type (
	BasicScope struct {
		scopes []map[string]eval.PValue
	}

	parentedScope struct {
		BasicScope
		parent eval.Scope
	}
)

// NewScope creates a new Scope instance that in turn consists of a stack of ephemeral scopes
func NewScope() eval.Scope {
	return &BasicScope{[]map[string]eval.PValue{make(map[string]eval.PValue, 8)}}
}

// NewParentedScope creates a scope that will override its parent. When a value isn't found in this
// scope, the search continues in the parent scope.
//
// All new or updated values will end up in this scope, i.e. no modifications are ever propagated to
// the parent scope.
func NewParentedScope(parent eval.Scope) eval.Scope {
	return &parentedScope{BasicScope{[]map[string]eval.PValue{make(map[string]eval.PValue, 8)}}, parent}
}

func NewScope2(h *types.HashValue) eval.Scope {
	top := make(map[string]eval.PValue, h.Len())
	h.EachPair(func(k, v eval.PValue) { top[k.String()] = v })
	return &BasicScope{[]map[string]eval.PValue{top}}
}

// No key can ever start with '::' or a capital letter
var groupKey = `::R`

func (e *BasicScope) RxGet(index int) (value eval.PValue, found bool) {
	// Variable is in integer form. An attempt is made to find a Regexp result group
	// in this scope using the special key '::R'. No attempt is made to traverse parent
	// scopes.
	if r, ok := e.scopes[len(e.scopes)-1][groupKey]; ok {
		if gv, ok := r.(*types.ArrayValue); ok && index < gv.Len() {
			return gv.At(index), true
		}
	}
	return eval.UNDEF, false
}

func (e *BasicScope) WithLocalScope(producer eval.Producer) eval.PValue {
	epCount := len(e.scopes)
	defer func() {
		// Pop all ephemerals
		e.scopes = e.scopes[:epCount]
	}()
	e.scopes = append(e.scopes, make(map[string]eval.PValue, 8))
	result := producer()
	return result
}

func (e *BasicScope) Get(name string) (value eval.PValue, found bool) {
	if strings.HasPrefix(name, `::`) {
		if value, found = e.scopes[0][name[2:]]; found {
			return
		}
		return eval.UNDEF, false
	}

	for idx := len(e.scopes) - 1; idx >= 0; idx-- {
		if value, found = e.scopes[idx][name]; found {
			return
		}
	}
	return eval.UNDEF, false
}

func (e *BasicScope) RxSet(variables []string) {
	// Assign the regular expression groups to an array value using the special key
	// '::R'. This overwrites an previous assignment in this scope
	varStrings := make([]eval.PValue, len(variables))
	for idx, v := range variables {
		varStrings[idx] = types.WrapString(v)
	}
	e.scopes[len(e.scopes)-1][groupKey] = types.WrapArray(varStrings)
}

func (e *BasicScope) Set(name string, value eval.PValue) bool {
	var current map[string]eval.PValue
	if strings.HasPrefix(name, `::`) {
		name = name[2:]
		current = e.scopes[0]
	} else {
		current = e.scopes[len(e.scopes)-1]
	}
	if _, found := current[name]; !found {
		current[name] = value
		return true
	}
	return false
}

func (e *parentedScope) Get(name string) (value eval.PValue, found bool) {
	value, found = e.BasicScope.Get(name)
	if !found {
		value, found = e.parent.Get(name)
	}
	return
}
