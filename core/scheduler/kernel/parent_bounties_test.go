package kernel_test

import (
	"github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/core/scheduler/kernel"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

func (s *KernelSuite) TestParentBountyNew(c *C) {
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
	}, {
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 2, Deps: []int64{1}},
			{Bounty: 3, Deps: []int64{2}},
			{Bounty: 4, Deps: []int64{2}},
		},
		expect: map[int64]kernel.Node{
			1: {ParentBounties: 9},
			2: {ParentBounties: 7},
			3: {ParentBounties: 0},
			4: {ParentBounties: 0},
		},
	}, {
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 2, Deps: []int64{1}},
			{Bounty: 3, Deps: []int64{1}},
			{Bounty: 4, Deps: []int64{2}},
		},
		expect: map[int64]kernel.Node{
			1: {ParentBounties: 9},
			2: {ParentBounties: 4},
			3: {ParentBounties: 0},
			4: {ParentBounties: 0},
		},
	}, {
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 2, Deps: []int64{1}},
			{Bounty: 3, Deps: []int64{1}},
			{Bounty: 4, Deps: []int64{2}},
		},
		expect: map[int64]kernel.Node{
			1: {ParentBounties: 9},
			2: {ParentBounties: 4},
			3: {ParentBounties: 0},
			4: {ParentBounties: 0},
		},
	}, {
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 1},
			{Bounty: 1},
			{Bounty: 1},
			{Bounty: 1, Deps: []int64{1}},
			{Bounty: 1, Deps: []int64{1, 2}},
			{Bounty: 1, Deps: []int64{2, 3}},
			{Bounty: 1, Deps: []int64{4}},
			{Bounty: 1, Deps: []int64{4, 5}},
			{Bounty: 1, Deps: []int64{2}},
			{Bounty: 1, Deps: []int64{2, 3}},
		},
		expect: map[int64]kernel.Node{
			1:  {ParentBounties: 3},
			2:  {ParentBounties: 4},
			3:  {ParentBounties: 2},
			4:  {ParentBounties: 2},
			5:  {ParentBounties: 1},
			6:  {ParentBounties: 0},
			7:  {ParentBounties: 0},
			8:  {ParentBounties: 0},
			9:  {ParentBounties: 0},
			10: {ParentBounties: 0},
			11: {ParentBounties: 0},
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

func (s *KernelSuite) TestParentBountyExisting(c *C) {
	for i, t := range []struct {
		given  []message.Intention
		expect map[int64]kernel.Node
	}{{
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 2, Deps: []int64{1}},
			{Bounty: 3, Deps: []int64{1}},
			{ID: 2, Deps: []int64{1}, Bounty: 5},
		},
		expect: map[int64]kernel.Node{
			1: {ParentBounties: 8},
			2: {ParentBounties: 0},
			3: {ParentBounties: 0},
		},
	}, {
		given: []message.Intention{
			{Bounty: 1},
			{Bounty: 2, Deps: []int64{1}},
			{Bounty: 3, Deps: []int64{2}},
			{ID: 2, Bounty: 2},
		},
		expect: map[int64]kernel.Node{
			1: {ParentBounties: 0},
			2: {ParentBounties: 3},
			3: {ParentBounties: 0},
		},
		// }, {
		// 	given: []message.Intention{
		// 		{Bounty: 1},
		// 		{Bounty: 2, Deps: []int64{1}},
		// 		{Bounty: 3, Deps: []int64{2}},
		// 		{Bounty: 4, Deps: []int64{2}},
		// 	},
		// 	expect: map[int64]kernel.Node{
		// 		1: {ParentBounties: 9},
		// 		2: {ParentBounties: 7},
		// 		3: {ParentBounties: 0},
		// 		4: {ParentBounties: 0},
		// 	},
		// }, {
		// 	given: []message.Intention{
		// 		{Bounty: 1},
		// 		{Bounty: 2, Deps: []int64{1}},
		// 		{Bounty: 3, Deps: []int64{1}},
		// 		{Bounty: 4, Deps: []int64{2}},
		// 	},
		// 	expect: map[int64]kernel.Node{
		// 		1: {ParentBounties: 9},
		// 		2: {ParentBounties: 4},
		// 		3: {ParentBounties: 0},
		// 		4: {ParentBounties: 0},
		// 	},
		// }, {
		// 	given: []message.Intention{
		// 		{Bounty: 1},
		// 		{Bounty: 2, Deps: []int64{1}},
		// 		{Bounty: 3, Deps: []int64{1}},
		// 		{Bounty: 4, Deps: []int64{2}},
		// 	},
		// 	expect: map[int64]kernel.Node{
		// 		1: {ParentBounties: 9},
		// 		2: {ParentBounties: 4},
		// 		3: {ParentBounties: 0},
		// 		4: {ParentBounties: 0},
		// 	},
		// }, {
		// 	given: []message.Intention{
		// 		{Bounty: 1},
		// 		{Bounty: 1},
		// 		{Bounty: 1},
		// 		{Bounty: 1},
		// 		{Bounty: 1, Deps: []int64{1}},
		// 		{Bounty: 1, Deps: []int64{1, 2}},
		// 		{Bounty: 1, Deps: []int64{2, 3}},
		// 		{Bounty: 1, Deps: []int64{4}},
		// 		{Bounty: 1, Deps: []int64{4, 5}},
		// 		{Bounty: 1, Deps: []int64{2}},
		// 		{Bounty: 1, Deps: []int64{2, 3}},
		// 	},
		// 	expect: map[int64]kernel.Node{
		// 		1:  {ParentBounties: 3},
		// 		2:  {ParentBounties: 4},
		// 		3:  {ParentBounties: 2},
		// 		4:  {ParentBounties: 2},
		// 		5:  {ParentBounties: 1},
		// 		6:  {ParentBounties: 0},
		// 		7:  {ParentBounties: 0},
		// 		8:  {ParentBounties: 0},
		// 		9:  {ParentBounties: 0},
		// 		10: {ParentBounties: 0},
		// 		11: {ParentBounties: 0},
		// 	},
	}} {
		c.Logf("test %d", i)

		k := kernel.New()

		for _, in := range t.given {
			k.Add(in)
		}

		c.Check(kernel.Nodes(k), jc.DeepEquals, t.expect)
	}
}
