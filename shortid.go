// Copyright 2016 Ventu.io. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

package shortid

import (
	randc "crypto/rand"
	"errors"
	"fmt"
	"math"
	randm "math/rand"
	"sync"
	"time"
)

type (
	Abc struct {
		sync.Mutex
		alphabet []rune
		seed     uint64
		original *Abc
	}

	Shortid struct {
		mx      sync.Mutex
		abc     *Abc
		worker  uint
		version uint      // restarts every year, rotates every 16 years!
		epoch   time.Time // restarts every year
		ts      uint      // timestamp (arbitrary units) incrementing since startts
		count   uint      // request count within ts
	}
)

const DEFAULT_ABC = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

func New(worker uint8, alphabet string, seed uint64) (*Shortid, error) {
	if abc, err := NewAbc(alphabet, seed); err == nil {
		now := time.Now()
		return &Shortid{
			abc:     abc,
			worker:  uint(worker),
			ts:      0,
			count:   0,
			epoch:   time.Date(now.Year()-1, time.December, 1, 1, 1, 1, 1, time.Local),
			version: uint(now.Year() % 16),
		}, nil
	} else {
		return nil, err
	}
}

func (sid *Shortid) Generate() string {
	return sid.generate(time.Now(), sid.epoch, sid.version)
}

func (sid *Shortid) generate(tm time.Time, epoch time.Time, version uint) string {
	ts := uint(tm.Sub(epoch).Seconds())
	var countrunes []rune
	sid.mx.Lock()
	if ts == sid.ts {
		sid.count++
	} else {
		sid.count = 0
		sid.ts = ts
	}
	if sid.count > 0 {
		countrunes = sid.abc.Encode(sid.count, 6)
	}
	sid.mx.Unlock()

	res := sid.abc.Encode(version, 4)
	res = append(res, sid.abc.Encode(sid.worker, 5)...)
	res = append(res, countrunes...)
	res = append(res, sid.abc.Encode(ts, 4)...)
	return string(res)
}

func (sid *Shortid) Abc() *Abc {
	return sid.abc
}

func NewAbc(alphabet string, seed uint64) (*Abc, error) {
	runes := []rune(alphabet)
	if len(runes) != len(DEFAULT_ABC) {
		return nil, errors.New(fmt.Sprintf("alphabet must contain %v unique characters", len(DEFAULT_ABC)))
	}
	if !areUnique(runes) {
		return nil, errors.New("alphabet must contain unique characters only")
	}
	original := &Abc{alphabet: runes, seed: seed}
	res := Abc{alphabet: nil, seed: seed, original: original}
	res.shuffle()
	return &res, nil
}

func areUnique(runes []rune) bool {
	found := make(map[rune]struct{})
	for _, r := range runes {
		if _, seen := found[r]; !seen {
			found[r] = struct{}{}
		}
	}
	return len(found) == len(runes)
}

func (abc *Abc) shuffle() {
	source := make([]rune, len(DEFAULT_ABC))
	copy(source, abc.original.alphabet)
	// abc.next(len(source)) // copied from the original code, useless?
	for len(source) > 1 {
		i := abc.next(len(source))
		abc.alphabet = append(abc.alphabet, source[i])
		source = append(source[:i], source[i+1:]...)
	}
	abc.alphabet = append(abc.alphabet, source[0])
}

func (abc *Abc) Reset() {
	abc.Lock()
	abc.seed = abc.original.seed
	abc.shuffle()
	abc.Unlock()
}

// Based on The Central Randomizer 1.3 (C) 1997 by Paul Houle (houle@msc.cornell.edu)
func (abc *Abc) next(lessthen int) int {
	abc.seed = (abc.seed*9301 + 49297) % 233280
	return int(math.Floor(float64(abc.seed) / (233280.0) * float64(lessthen)))
}

func (abc *Abc) Encode(val uint, dig uint) []rune {
	if dig > 6 {
		panic("max dig=6")
	}
	size := 1
	if val >= 1 {
		size = int(math.Log2(float64(val)))/int(dig) + 1
		if size < 1 {
			size = 1
		}
	}

	mask := int(1<<dig - 1)
	buf := make([]byte, size)
	if _, err := randc.Read(buf); err != nil {
		for i, _ := range buf {
			buf[i] = byte(randm.Int31n(0xff))
		}
	}
	res := make([]rune, size)
	abc.Lock()
	for i, _ := range res {
		shift := dig * uint(i)
		index := int(val>>shift) & mask
		if dig < 6 {
			index = index | (int(buf[i]) & int(0x3f-mask))
		}
		res[i] = abc.alphabet[index]
	}
	abc.Unlock()
	return res
}
