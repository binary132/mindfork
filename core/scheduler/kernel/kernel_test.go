package kernel_test

import (
	"errors"

	"github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/core/scheduler"
	"github.com/mindfork/mindfork/core/scheduler/kernel"
	mfm "github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = scheduler.Scheduler(&kernel.Kernel{})

//TODO: test Free
//TODO: revise / test Orderings

func (s *KernelSuite) TestParentBounty(c *C) {
	for i, t := range []struct {
		given  []message.Intention
		expect map[int64]kernel.Node
	}{{
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 2, Deps: []int64{1}},
		},
		expect: map[int64]kernel.Node{
			1: {ParentBounties: 2},
			2: {ParentBounties: 0},
		},
	}, {
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 2, Deps: []int64{1}},
			{Bounty: 3, Deps: []int64{1}},
		},
		expect: map[int64]kernel.Node{
			1: {ParentBounties: 5},
			2: {ParentBounties: 0},
			3: {ParentBounties: 0},
		},
	}, {
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 2, Deps: []int64{1}},
			{Bounty: 3, Deps: []int64{2}},
		},
		expect: map[int64]kernel.Node{
			1: {ParentBounties: 5},
			2: {ParentBounties: 3},
			3: {ParentBounties: 0},
		},
	}} {
		c.Logf("test %d", i)

		k := kernel.New()

		for _, in := range t.given {
			k.Add(in)
		}

		c.Check(kernel.Nodes(k), jc.DeepEquals, t.expect)
	}
}

func (s *KernelSuite) TestRoots(c *C) {
	for i, t := range []struct {
		should string
		given  []message.Intention
		expect map[int64]message.Intention
	}{{
		should: "make a few new roots",
		given:  []message.Intention{{}, {}, {}},
		expect: map[int64]message.Intention{
			1: message.Intention{ID: 1},
			2: message.Intention{ID: 2},
			3: message.Intention{ID: 3},
		},
	}, {
		should: "make and remove a few roots",
		given:  []message.Intention{{}, {}, {ID: 2, Deps: []int64{1}}},
		expect: map[int64]message.Intention{
			2: message.Intention{ID: 2, Deps: []int64{1}},
		},
	}, {
		should: "show only one root for a tree",
		given: []message.Intention{
			{},
			{Deps: []int64{1}},
			{Deps: []int64{2}},
			{Deps: []int64{3}},
		},
		expect: map[int64]message.Intention{
			4: message.Intention{ID: 4, Deps: []int64{3}},
		},
	}, {
		should: "make an orphaned child of two parents into a root",
		given: []message.Intention{
			{}, {}, {},
			{ID: 1, Deps: []int64{3}},
			{ID: 2, Deps: []int64{3}},
			{ID: 1},
			{ID: 2},
		},
		expect: map[int64]message.Intention{
			1: message.Intention{ID: 1},
			2: message.Intention{ID: 2},
			3: message.Intention{ID: 3},
		},
	}, {
		should: "make two orphaned children of a parent into roots",
		given: []message.Intention{
			{}, {}, {}, {},
			{ID: 1, Deps: []int64{3, 4}},
			{ID: 2, Deps: []int64{3, 4}},
			{ID: 1},
			{ID: 2},
		},
		expect: map[int64]message.Intention{
			1: message.Intention{ID: 1},
			2: message.Intention{ID: 2},
			3: message.Intention{ID: 3},
			4: message.Intention{ID: 4},
		},
	}, {
		should: "make and split a tree into two",
		given: []message.Intention{
			{}, {}, {},
			{Deps: []int64{1, 2, 3}},
			{Deps: []int64{4}},
			{ID: 5, Deps: []int64{1, 2, 3}},
		},
		expect: map[int64]message.Intention{
			4: message.Intention{ID: 4, Deps: []int64{1, 2, 3}},
			5: message.Intention{ID: 5, Deps: []int64{1, 2, 3}},
		},
	}, {
		should: "make a tree, split it into two, and rejoin it",
		given: []message.Intention{
			{}, {},
			{Deps: []int64{1, 2}},
			{Deps: []int64{3}},
			{},
			{ID: 4, Deps: []int64{1, 2}},
			{ID: 3, Deps: []int64{4, 5}},
		},
		expect: map[int64]message.Intention{
			3: message.Intention{ID: 3, Deps: []int64{4, 5}},
		},
	}} {
		c.Logf("test %d: should %s", i, t.should)

		k := kernel.New()

		for _, i := range t.given {
			k.Add(i)
		}

		c.Check(kernel.Roots(k), jc.DeepEquals, t.expect)
	}
}

func (cs *KernelSuite) TestAdd(c *C) {
	for i, t := range []struct {
		should string
		given  []message.Intention
		expect []mfm.Message
	}{{
		should: "fail to Add an Intention whose ID does not exist",
		given:  []message.Intention{{ID: 1}},
		expect: []mfm.Message{message.Error{
			Err: errors.New("no such Intention 1"),
		}},
	}, {
		should: "fail to Add an Intention with invalid Deps",
		given:  []message.Intention{{Deps: []int64{1}}},
		expect: []mfm.Message{message.Error{
			Err: errors.New("no such Intention 1"),
		}},
	}, {
		should: "Add a single new Intention",
		given:  []message.Intention{{}},
		expect: []mfm.Message{message.Intention{ID: 1}},
	}, {
		should: "Add a few new Intentions",
		given:  []message.Intention{{}, {}, {}},
		expect: []mfm.Message{
			message.Intention{ID: 1},
			message.Intention{ID: 2},
			message.Intention{ID: 3},
		},
	}, {
		should: "fail to Add an Intention with some invalid Deps",
		given: []message.Intention{
			{},
			{Deps: []int64{1}},
			{Deps: []int64{2}},
			{Deps: []int64{3, 5}},
			{Deps: []int64{2, 3}},
		},
		expect: []mfm.Message{
			message.Intention{ID: 1},
			message.Intention{ID: 2, Deps: []int64{1}},
			message.Intention{ID: 3, Deps: []int64{2}},
			message.Error{Err: errors.New("no such Intention 5")},
			message.Intention{ID: 4, Deps: []int64{2, 3}},
		},
	}, {
		should: "fail to Add an Intention with a dep cycle",
		given: []message.Intention{
			{},
			{Deps: []int64{1}},
			{Deps: []int64{2}},
			{ID: 1, Deps: []int64{3}},
		},
		expect: []mfm.Message{
			message.Intention{ID: 1},
			message.Intention{ID: 2, Deps: []int64{1}},
			message.Intention{ID: 3, Deps: []int64{2}},
			message.Error{Err: errors.New("cycle requested: [2 1]")},
		},
	}} {
		c.Logf("test %d: should %s\n  given: %v\n  expect: %v",
			i, t.should, t.given, t.expect,
		)

		k := kernel.New()

		result := make([]mfm.Message, len(t.given))

		for i, in := range t.given {
			result[i] = k.Add(in)
		}

		c.Check(result, jc.DeepEquals, t.expect)
	}
}

func (s *KernelSuite) TestAvailable(c *C) {
	for i, t := range []struct {
		should        string
		given         []message.Intention
		givenOrdering scheduler.Ordering
		expect        []message.Intention
	}{{
		should: "show a few free Intentions",
		given:  []message.Intention{{}, {}, {}},
		expect: []message.Intention{
			message.Intention{ID: 1},
			message.Intention{ID: 2},
			message.Intention{ID: 3},
		},
	}, {
		should: "not show non-free Intentions",
		given:  []message.Intention{{}, {}, {Deps: []int64{1}}},
		expect: []message.Intention{
			message.Intention{ID: 1},
			message.Intention{ID: 2},
		},
	}, {
		should: "not show non-free Intentions after mutation",
		given:  []message.Intention{{}, {}, {ID: 2, Deps: []int64{1}}},
		expect: []message.Intention{
			message.Intention{ID: 1},
		},
	}, {
		should: "show re-freed Intentions after mutation",
		given: []message.Intention{
			{}, {},
			{ID: 2, Deps: []int64{1}},
			{ID: 2},
		},
		expect: []message.Intention{
			message.Intention{ID: 1},
			message.Intention{ID: 2},
		},
	}} {
		c.Logf("test %d: should %s", i, t.should)

		k := kernel.New()

		for _, i := range t.given {
			k.Add(i)
		}

		c.Check(k.Available(t.givenOrdering), jc.DeepEquals, t.expect)
	}
}

func (s *KernelSuite) TestCheckCycle(c *C) {
	type edge struct {
		from int64
		to   int64
	}

	for i, t := range []struct {
		should    string
		givenMap  map[int64]message.Intention
		givenEdge edge
		expect    []int64
	}{{
		should: "result in nil for no cycles",
		givenMap: map[int64]message.Intention{
			1: {ID: 1}, 2: {ID: 2},
		},
		givenEdge: edge{1, 2},
	}, {
		should: "result in nil for acyclical deps",
		givenMap: map[int64]message.Intention{
			1: {ID: 1, Deps: []int64{2}},
			2: {ID: 2},
			3: {ID: 3},
		},
		givenEdge: edge{1, 3},
	}, {
		should: "show cycle for simple reversed edge",
		givenMap: map[int64]message.Intention{
			1: {ID: 1, Deps: []int64{2}},
			2: {ID: 2},
		},
		givenEdge: edge{2, 1},
		expect:    []int64{1, 2},
	}, {
		should: "show cycle for triangular cycle",
		givenMap: map[int64]message.Intention{
			1: {ID: 1},
			2: {ID: 2, Deps: []int64{1}},
			3: {ID: 3, Deps: []int64{2}},
		},
		givenEdge: edge{1, 3},
		expect:    []int64{2, 1},
	}, {
		should: "show cycle for rectangular cycle",
		givenMap: map[int64]message.Intention{
			1: {ID: 1},
			2: {ID: 2, Deps: []int64{1}},
			3: {ID: 3, Deps: []int64{2}},
			4: {ID: 4, Deps: []int64{3}},
		},
		givenEdge: edge{1, 4},
		expect:    []int64{2, 1},
	}, {
		should: "show cycle for diamond cycle",
		givenMap: map[int64]message.Intention{
			1: {ID: 1, Deps: []int64{2, 3}},
			2: {ID: 2, Deps: []int64{4}},
			3: {ID: 3, Deps: []int64{4}},
			4: {ID: 4},
		},
		givenEdge: edge{4, 1},
		expect:    []int64{2, 4},
	}, {
		should: "show cycle for more complicated cycle",
		givenMap: map[int64]message.Intention{
			1: {ID: 1, Deps: []int64{2, 3, 4}},
			2: {ID: 2, Deps: []int64{4, 6}},
			3: {ID: 3, Deps: []int64{4}},
			4: {ID: 4, Deps: []int64{5}},
			5: {ID: 5},
			6: {ID: 6},
		},
		givenEdge: edge{5, 1},
		expect:    []int64{4, 5},
	}, {
		should: "return nil for no cycle in more complicated graph",
		givenMap: map[int64]message.Intention{
			1: {ID: 1, Deps: []int64{2, 3, 4}},
			2: {ID: 2, Deps: []int64{4, 6}},
			3: {ID: 3, Deps: []int64{4}},
			4: {ID: 4, Deps: []int64{5}},
			5: {ID: 5},
			6: {ID: 6},
		},
		givenEdge: edge{6, 4},
	}} {
		c.Logf("test %d: should %s", i, t.should)
		got := kernel.CheckCycle(
			t.givenMap,
			t.givenEdge.from,
			t.givenEdge.to,
		)
		c.Check(got, jc.DeepEquals, t.expect)
	}
}
