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
		ms      uint      // timestamp (arbitrary units) incrementing since startts
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
			ms:      0,
			count:   0,
			epoch:   time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local),
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
	ms := uint(tm.Sub(epoch).Nanoseconds() / 1000000)
	var countrunes []rune
	sid.mx.Lock()
	if ms == sid.ms {
		sid.count++
	} else {
		sid.count = 0
		sid.ms = ms
	}
	if sid.count > 0 {
		countrunes = sid.abc.Encode(sid.count, 0, 6)
	}
	sid.mx.Unlock()

	res := sid.abc.Encode(ms, 8, 5)
	res = append(res, sid.abc.Encode(sid.worker, 1, 5)[0])
	res = append(res, countrunes...)
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

func (abc *Abc) Encode(val uint, size uint, dig uint) []rune {
	if dig > 6 {
		panic("max dig=6")
	}

	var csize uint = 1
	if val >= 1 {
		csize = uint(math.Log2(float64(val))/float64(dig) + 1.0)
		if csize < 1 {
			csize = 1
		}
	}
	if size == 0 {
		size = csize
	} else if size < csize {
		panic("cannot accommodate data")
	}

	mask := int(1<<dig - 1)
	buf := make([]byte, int(size))
	if _, err := randc.Read(buf); err != nil {
		for i, _ := range buf {
			buf[i] = byte(randm.Int31n(0xff))
		}
	}
	res := make([]rune, int(size))
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
