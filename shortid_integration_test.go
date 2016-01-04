// Copyright (c) 2016 Ventu.io, Oleg Sklyar, contributors
// The use of this source code is governed by a MIT style license found in the LICENSE file

package shortid_test

import (
	"github.com/ventu-io/go-shortid"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestShortid_Generate_oneIdPerYear_over34Years_uniqueOf9Symbols(t *testing.T) {
	now := time.Now()
	sid := shortid.MustNew(0, shortid.DEFAULT_ABC, 155000)
	for years := 0; years < 38; years++ {
		tm := now.AddDate(years, 0, 1)
		id, err := sid.GenerateInternal(&tm, now)
		if years <= 34 && err != nil {
			t.Errorf("no error expected for lifespan %vy", years)
		} else if years > 34 && err == nil {
			t.Errorf("error expected for lifespan %vy", years)
		}
		if err == nil && len(id) != 9 {
			t.Errorf("all ids are expected to be of length 9, found %v for %vy", id, years)
		}
	}
}

func TestShortid_Generate_approx3MilIdsWith5MinStep_over33Years_unique9_11Symbols(t *testing.T) {
	// t.Skip("do not run for unit testing")

	sid := shortid.MustNew(0, shortid.DEFAULT_ABC, 155000)
	ids := make(map[string]struct{}, 3800000)
	now := time.Now()
	tm := now
	end := now.AddDate(34, 0, 1)
	maxlen := 0.
	minlen := 1e9
	for tm.Before(end) {
		// step: any value between 1ns and 10min+1ns (5min on average)
		tm = tm.Add(time.Duration(1 + rand.Int63n(600000000000)))
		id, err := sid.GenerateInternal(&tm, now)
		if err != nil {
			t.Errorf("error for %v: %v", tm, err)
		}
		if _, ok := ids[id]; !ok {
			ids[id] = struct{}{}
			maxlen = math.Max(maxlen, float64(len(id)))
			minlen = math.Min(minlen, float64(len(id)))
		} else {
			t.Errorf("duplicate on %v: %v", tm, id)
		}
	}
	if minlen != 9 {
		t.Errorf("min length expected to be 9, found %v", minlen)
	}
	if maxlen > 11 {
		t.Errorf("max length expected to be 11, found %v", maxlen)
	}
	t.Logf("generated %v Ids from %v till %v", len(ids), now, tm)
}

func TestShortid_Generate_500kValuesEach_at6Timepoints_unique9_11Symbols(t *testing.T) {
	// t.Skip("do not run for unit testing")

	var n int = 500000
	var m int = 6
	sid, _ := shortid.New(0, shortid.DEFAULT_ABC, 155000)
	ids := make(map[string]struct{}, n)
	tms := make([]time.Time, m)
	maxlen := 0.
	minlen := 1e9
	for j := 0; j < m; j++ {
		// heck knows when, but at some point in the future
		tms[j] = time.Now().Add(time.Duration(-rand.Int63n(60000000000000))).Add(time.Duration(rand.Int63n(6000)-9000) * time.Hour)
		for i := 0; i < n; i++ {
			id, _ := sid.GenerateInternal(nil, tms[j])
			if _, ok := ids[id]; !ok {
				ids[id] = struct{}{}
			}
			maxlen = math.Max(maxlen, float64(len(id)))
			minlen = math.Min(minlen, float64(len(id)))
		}
	}
	if len(ids) != n*m {
		t.Errorf("expected len 1e6, found %v. duplicates: %v", len(ids))
	}
	if minlen != 9 {
		t.Errorf("min length expected to be 9, found %v", minlen)
	}
	if maxlen < 10 || 12 < maxlen {
		t.Errorf("max length expected to be between 10 and 12, found %v", maxlen)
	}
	t.Logf("generated %v Ids with epochs at: %v", len(ids), tms)
}

func TestShortid_Generate_500kValues_concurrently(t *testing.T) {
	sid := shortid.MustNew(0, shortid.DEFAULT_ABC, 155000)
	ids := make(map[string]struct{}, 900000)
	var mx sync.Mutex
	generate := func(done chan bool) {
		for i := 0; i < 300000; i++ {
			id := sid.MustGenerate()
			mx.Lock()
			ids[id] = struct{}{}
			mx.Unlock()
		}
		done <- true
	}
	done1 := make(chan bool)
	go generate(done1)
	done2 := make(chan bool)
	go generate(done2)
	done3 := make(chan bool)
	go generate(done3)
	<-done1
	<-done2
	<-done3
	if len(ids) != 900000 {
		t.Errorf("expected %v unique ids, found %v", 900000, len(ids))
	}
}
