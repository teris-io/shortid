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
	t.Log(sid.generate(now, now.AddDate(0, 0, -1), 5))
	t.Log(sid.generate(now.AddDate(1, 0, 1), now, 5))
	t.Log(sid.generate(now.AddDate(15, 0, 1), now, 5))
	t.Log(sid.generate(now.AddDate(16, 0, 1), now, 5))
	t.Log(sid.generate(now.AddDate(34, 0, 1), now, 5))
}
