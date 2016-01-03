// Copyright 2016 Ventu.io. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

package shortid

import (
	"testing"
	"time"
)

func TestShortid_generatex4(t *testing.T) {
	sid, _ := New(0, DEFAULT_ABC, 155000)
	now := time.Now()
	id00 := sid.generate(now.AddDate(0, 0, 1), now, 5)
	id34 := sid.generate(now.AddDate(34, 0, 1), now, 5)
	if len(id00) != 9 || len(id34) != 9 {
		t.Errorf("all ids with different ms are expected to be of length 9")
	}
}
