// Copyright 2016 Ventu.io. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

// Seed computation based on The Central Randomizer 1.3
// (C) 1997 by Paul Houle (houle@msc.cornell.edu)

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
		alphabet []rune
		seed     uint64
	}

	Shortid struct {
		abc    *Abc
		worker uint
		epoch  time.Time  // ids can be generated for 34 years since this date
		ms     uint       // ms since epoch for the last id
		count  uint       // request count within the same ms
		mx     sync.Mutex // locks access to ms and count
	}
)

const DEFAULT_ABC = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

// New constructs a new instance of short id generator.
func New(worker uint8, alphabet string, seed uint64) (*Shortid, error) {
	if abc, err := NewAbc(alphabet, seed); err == nil {
		return &Shortid{
			abc:    abc,
			worker: uint(worker),
			epoch:  time.Date(2016, time.January, 1, 0, 0, 0, 0, time.Local),
			ms:     0,
			count:  0,
		}, nil
	} else {
		return nil, err
	}
}

func MustNew(worker uint8, alphabet string, seed uint64) *Shortid {
	if res, err := New(worker, alphabet, seed); err == nil {
		return res
	} else {
		panic(err)
	}
}

func (sid *Shortid) Generate() (string, error) {
	return sid.generate(time.Now(), sid.epoch)
}

func (sid *Shortid) generate(tm time.Time, epoch time.Time) (string, error) {
	ms, count := sid.msAndCountLocking(tm, epoch)
	var res []rune
	if tmp, err := sid.abc.Encode(ms, 8, 5); err != nil {
		return "", err
	} else {
		res = tmp
	}
	if tmp, err := sid.abc.Encode(sid.worker, 1, 5); err != nil {
		return "", err
	} else {
		res = append(res, tmp[0])
	}
	if count > 0 {
		if countrunes, err := sid.abc.Encode(count, 0, 6); err == nil {
			res = append(res, countrunes...)
		} else {
			return "", err
		}
	}
	return string(res), nil
}

func (sid *Shortid) msAndCountLocking(tm time.Time, epoch time.Time) (uint, uint) {
	sid.mx.Lock()
	defer sid.mx.Unlock()
	ms := uint(tm.Sub(epoch).Nanoseconds() / 1000000)
	if ms == sid.ms {
		sid.count++
	} else {
		sid.count = 0
		sid.ms = ms
	}
	return sid.ms, sid.count
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
	res := Abc{alphabet: nil, seed: seed}
	res.shuffle(alphabet)
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

func (abc *Abc) shuffle(alphabet string) {
	source := []rune(alphabet)
	for len(source) > 1 {
		abc.seed = (abc.seed*9301 + 49297) % 233280
		i := int(abc.seed * uint64(len(source)) / 233280)

		abc.alphabet = append(abc.alphabet, source[i])
		source = append(source[:i], source[i+1:]...)
	}
	abc.alphabet = append(abc.alphabet, source[0])
}

func (abc *Abc) Encode(val, size, digits uint) ([]rune, error) {
	if digits < 3 || 6 < digits {
		return nil, errors.New(fmt.Sprintf("allowed digits range [3,6], found %v", digits))
	}

	var computedSize uint = 1
	if val >= 1 {
		computedSize = uint(math.Log2(float64(val)))/digits + 1
	}
	if size == 0 {
		size = computedSize
	} else if size < computedSize {
		return nil, errors.New(fmt.Sprintf("cannot accommodate data, need %v digits, got %v", computedSize, size))
	}

	buf := make([]byte, int(size))
	if _, err := randc.Read(buf); err != nil {
		for i, _ := range buf {
			buf[i] = byte(randm.Int31n(0xff))
		}
	}

	mask := int(1<<digits - 1)
	res := make([]rune, int(size))
	for i, _ := range res {
		shift := digits * uint(i)
		index := int(val>>shift) & mask
		if digits < 6 {
			index = index | (int(buf[i]) & int(0x3f-mask))
		}
		res[i] = abc.alphabet[index]
	}
	return res, nil
}

func (abc *Abc) MustEncode(val, size, digits uint) []rune {
	if res, err := abc.Encode(val, size, digits); err == nil {
		return res
	} else {
		panic(err)
	}
}
