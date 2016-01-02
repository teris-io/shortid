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
		mx     sync.Mutex
		abc    *Abc
		worker uint
		sec    uint
		count  uint
	}
)

const DEFAULT_ABC = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

var (
	version  uint      // restarts every 16 years
	secstart time.Time // restarts every year
)

func init() {
	now := time.Now()
	secstart = time.Date(now.Year()-1, time.December, 1, 1, 1, 1, 1, time.Local)
	version = uint(now.Year() % 16)
}

func New(worker uint8, alphabet string, seed uint64) (*Shortid, error) {
	if abc, err := NewAbc(alphabet, seed); err == nil {
		return &Shortid{abc: abc, worker: uint(worker), sec: 0, count: 0}, nil
	} else {
		return nil, err
	}
}

func (sid *Shortid) Generate() string {
	sec := uint(time.Now().Sub(secstart).Seconds() * 1e-3)

	var countrunes []rune
	sid.mx.Lock()
	if sec == sid.sec {
		sid.count++
	} else {
		sid.count = 0
		sid.sec = sec
	}
	if sid.count > 0 {
		countrunes = sid.abc.Encode(sid.count)
	}
	sid.mx.Unlock()

	res := sid.abc.Encode(version)
	res = append(res, sid.abc.Encode(sid.worker)...)
	res = append(res, countrunes...)
	res = append(res, sid.abc.Encode(sec)...)
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

func (abc *Abc) Encode(val uint) []rune {
	nlookups := int(math.Log2(float64(val))-5) - 1
	if nlookups < 1 {
		nlookups = 1
	}
	buf := make([]byte, nlookups)
	_, err := randc.Read(buf)
	if err != nil {
		for i, _ := range buf {
			buf[i] = byte(randm.Int31n(16))
		}
	}
	res := make([]rune, nlookups)
	abc.Lock()
	for i, _ := range res {
		p1 := uint8(val>>uint(4*i)) & 0x0f
		p2 := uint8(buf[i]) & 0x30
		res[i] = abc.alphabet[p1|p2]
	}
	abc.Unlock()
	return res
}
