// Copyright 2016 Ventu.io. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

package shortid_test

import (
	"github.com/ventu-io/go-shortid"
	"testing"
	"math"
)

func TestShortid_Generate1MilValues_unique(t *testing.T) {
	n := int(5e5)
	sid, _ := shortid.New(0, shortid.DEFAULT_ABC, 155000)
	ids := make(map[string]struct{})
	maxlen := 0.
	minlen := 1e9
	for i := 0; i < n; i++ {
		id := sid.Generate()
		if _, ok := ids[id]; !ok {
			ids[id] = struct{}{}
		}
		maxlen = math.Max(maxlen, float64(len(id)))
		minlen = math.Min(minlen, float64(len(id)))
	}
	if len(ids) != n {
		t.Errorf("expected len 1e6, found %v. duplicates: %v", len(ids))
	}
	if minlen != 9 {
		t.Errorf("min length expected to be 9, found %v", minlen)
	}
	if maxlen > 11 {
		t.Errorf("max length expected to be 11, found %v", maxlen)
	}
}
