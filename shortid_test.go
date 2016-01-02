// Copyright 2016 Ventu.io. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

package shortid_test

import (
	"github.com/ventu-io/go-shortid"
	"testing"
)

func TestShortid_Generate1MilValues_unique(t *testing.T) {
	n := int(1e6)
	sid, _ := shortid.New(0, shortid.DEFAULT_ABC, 155000)
	ids := make(map[string]struct{})
	var dups []string
	for i := 0; i < n; i++ {
		id := sid.Generate()
		if _, ok := ids[id]; !ok {
			ids[id] = struct{}{}
		} else {
			dups = append(dups, id)
		}
		if i <= 10 || i >= n-10 {
			t.Log(id)
		}
	}
	if len(ids) != n {
		t.Errorf("expected len 1e6, found %v. duplicates: %v", len(ids), dups)
	}
}
