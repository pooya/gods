// Copyright 2013 Shayan Pooya. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"sort"
	"testing"
)

const (
	count  = 30000
	delNum = 5000
)

func TestBasic1(t *testing.T) {
	list := make([]int, 0)
	treap := &Treap{}

	for i := 0; i < count; i++ {
		k := i
		addIt := true
		for _, j := range list {
			if j == k {
				addIt = false
				break
			}
		}
		if addIt {
			treap.Insert(k)
			list = append(list, k)
		}
	}
	uniqeElements := len(list)

	traverse := treap.Traverse1()
	if !sort.IntsAreSorted(traverse) {
		t.Error("not sorted", traverse)
	}

	traverse = treap.Traverse2()
	if !sort.IntsAreSorted(traverse) {
		t.Error("not sorted", traverse)
	}

	traverse = treap.Traverse3()
	if !sort.IntsAreSorted(traverse) {
		t.Error("not sorted", traverse)
	}

	index := 0
	for i := 0; i < delNum; i++ {
		key := list[index]
		if !sort.IntsAreSorted(treap.Traverse3()) {
			t.Error("Traverse3 is not sorted.")
		}
		removed := treap.Remove(key)
		if !removed {
			t.Error("Could not remove key: ", key, " with index ", index)
			return
		}
		index += uniqeElements / delNum

		if index >= uniqeElements {
			break
		}
	}
}
