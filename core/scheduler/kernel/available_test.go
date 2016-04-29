package kernel_test

import (
	"time"

	"github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/core/scheduler"
	"github.com/mindfork/mindfork/core/scheduler/kernel"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

func (s *KernelSuite) TestAvailable(c *C) {
	times := make([]*time.Time, 5)
	for i := range times {
		times[i] = new(time.Time)
		*times[i] = time.Now().Add(time.Duration(i) * time.Second)
	}

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
	}, {
		should: "sort available Intentions by ID",
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 2},
			{Bounty: 3},
			{Bounty: 4},
			{Bounty: 5, Deps: []int64{1}},
		},
		givenOrdering: kernel.ByID,
		expect: []message.Intention{
			message.Intention{ID: 1, Bounty: 1},
			message.Intention{ID: 2, Bounty: 2},
			message.Intention{ID: 3, Bounty: 3},
			message.Intention{ID: 4, Bounty: 4},
		},
	}, {
		should: "sort available Intentions by Bounty",
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 3},
			{Bounty: 2},
			{Bounty: 4},
			{Bounty: 5, Deps: []int64{1}},
		},
		givenOrdering: kernel.ByBounty,
		expect: []message.Intention{
			message.Intention{ID: 4, Bounty: 4},
			message.Intention{ID: 2, Bounty: 3},
			message.Intention{ID: 3, Bounty: 2},
			message.Intention{ID: 1, Bounty: 1},
		},
	}, {
		should: "sort available Intentions by date",
		given: []message.Intention{
			{When: times[0], Bounty: 1},
			{When: times[0], Bounty: 2},
			{When: times[2], Bounty: 2},
			{When: times[1], Bounty: 3},
			{When: times[3], Bounty: 4},
			{Bounty: 5},
			{Bounty: 6},
			{When: times[4], Bounty: 5, Deps: []int64{1}},
		},
		givenOrdering: kernel.ByDate,
		expect: []message.Intention{
			{ID: 2, When: times[0], Bounty: 2},
			{ID: 1, When: times[0], Bounty: 1},
			{ID: 4, When: times[1], Bounty: 3},
			{ID: 3, When: times[2], Bounty: 2},
			{ID: 5, When: times[3], Bounty: 4},
			{ID: 7, Bounty: 6},
			{ID: 6, Bounty: 5},
		},
	}, {
		should: "sort available Intentions by combined score",
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 1},
			{Bounty: 1},
			{Bounty: 1},
			{Bounty: 2, Deps: []int64{1, 2}},
			{Bounty: 3, Deps: []int64{3, 4}},
			{Bounty: 5, Deps: []int64{6}},
		},
		givenOrdering: kernel.ByScore,
		expect: []message.Intention{
			message.Intention{ID: 1, Bounty: 1},
			message.Intention{ID: 2, Bounty: 1},
			message.Intention{ID: 3, Bounty: 1},
			message.Intention{ID: 4, Bounty: 1},
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
