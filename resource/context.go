package resource

import (
	"github.com/puppetlabs/go-evaluator/eval"
	"github.com/puppetlabs/go-parser/issue"
	"gonum.org/v1/gonum/graph"
	"github.com/puppetlabs/go-evaluator/types"
)

const (
	SHARED_MAP = `sharedMap`
	NODE_GRAPH = `nodeGraph`
	NODE_JOBS = `nodeJobs`
	RESOURCES = `resources`
	CURRENT_NODE = `currentNode`
	EXTERNAL_EDGES_TO = `externalTo`
)

// GetGraph returns concurrent graph that is shared between all contexts
func GetGraph(c eval.Context) Graph {
	if g, ok := getSharedMap(c)[NODE_GRAPH]; ok {
		return g.(Graph)
	}
	panic(eval.Error(c, eval.EVAL_MISSING_REQUIRED_CONTEXT_VARIABLE, issue.H{`key`: NODE_GRAPH}))
}

// Resources returns the resources that has been created so far in the callers context. The
// hash is sorted by file name and location of where the resource was instantitated.
func Resources(c eval.Context) eval.KeyedValue {
	rs := getResources(c)
	entries := make([]*types.HashEntry, 0, len(rs))
	for k, r := range rs {
		entries = append(entries, types.WrapHashEntry2(k, r))
	}
	sortByEntriesLocation(entries)
	return types.WrapHash(entries)
}

func defineResource(c eval.Context, resource eval.PuppetObject, location issue.Location) {
	rs := getResources(c)
	ref := Reference(c, resource)
	if oh, ok := rs[ref]; ok {
		if oh.value != nil {
			panic(eval.Error(c, EVAL_DUPLICATE_RESOURCE, issue.H{`ref`: ref, `previous_location`: issue.LocationString(oh.location)}))
		}
		oh.value = resource
		oh.location = location
	} else {
		rs[ref] = &handle{resource, location}
	}
}

func getCurrentNode(c eval.Context) *node {
	if g, ok := c.Get(CURRENT_NODE); ok {
		return g.(*node)
	}
	panic(eval.Error(c, eval.EVAL_MISSING_REQUIRED_CONTEXT_VARIABLE, issue.H{`key`: CURRENT_NODE}))
}

func setCurrentNode(c eval.Context, n *node) {
	c.Set(CURRENT_NODE, n)
}

func getExternalEdgesTo(c eval.Context) []graph.Node {
	if rs, ok := c.Get(EXTERNAL_EDGES_TO); ok {
		return rs.([]graph.Node)
	}
	panic(eval.Error(c, eval.EVAL_MISSING_REQUIRED_CONTEXT_VARIABLE, issue.H{`key`: EXTERNAL_EDGES_TO}))
}

func setExternalEdgesTo(c eval.Context, edges []graph.Node) {
	c.Set(EXTERNAL_EDGES_TO, edges)
}

func getResources(c eval.Context) map[string]*handle {
	if rs, ok := c.Get(RESOURCES); ok {
		return rs.(map[string]*handle)
	}
	panic(eval.Error(c, eval.EVAL_MISSING_REQUIRED_CONTEXT_VARIABLE, issue.H{`key`: RESOURCES}))
}

func setResources(c eval.Context, resources map[string]*handle) {
	c.Set(RESOURCES, resources)
}

// Things shared
func getSharedMap(c eval.Context) map[string]interface{} {
	if g, ok := c.Get(SHARED_MAP); ok {
		return g.(map[string]interface{})
	}
	panic(eval.Error(c, eval.EVAL_MISSING_REQUIRED_CONTEXT_VARIABLE, issue.H{`key`: SHARED_MAP}))
}

func getNodeJobs(c eval.Context) chan *nodeJob {
	if g, ok := getSharedMap(c)[NODE_JOBS]; ok {
		return g.(chan *nodeJob)
	}
	panic(eval.Error(c, eval.EVAL_MISSING_REQUIRED_CONTEXT_VARIABLE, issue.H{`key`: NODE_JOBS}))
}

func sortByEntriesLocation(entries []*types.HashEntry) {
	v := make([]eval.PValue, len(entries))
	for i, e := range entries {
		v[i] = e
	}
	types.WrapArray(v).Sort(func(a, b eval.PValue) bool {
		l1 := a.(*types.HashEntry).Value().(issue.Located).Location()
		l2 := b.(*types.HashEntry).Value().(issue.Located).Location()
		if l1.File() == l2.File() {
			ld := l1.Line() - l2.Line()
			if ld == 0 {
				return l1.Pos() < l2.Pos()
			}
			return ld < 0
		}
		return l1.File() < l2.File()
	}).EachWithIndex(func(e eval.PValue, i int) {
		entries[i] = e.(*types.HashEntry)
	})
}
