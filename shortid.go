// Copyright 2016 Ventu.io. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

// Original idea of the algorithm, although reworked:
// Copyright (c) 2015 Dylan Greene, contributors
// https://github.com/dylang/shortid

// Seed computation based on The Central Randomizer 1.3
// (C) 1997 by Paul Houle (houle@msc.cornell.edu)

// Package shortid provides functionality to generate short (normally 9 to 11 symbols), unique (for
// 34 years from 1/1/2016) and URL friendly (by default) Ids.
//
// Being inspired by the node.js shortid library (https://github.com/dylang/shortid) this package is
// not a simple Go port, it actually adds a number of improvements to the original: (i) safe to be
// used in concurrent goroutines, (ii) no yearly epoch resets are required for 34 years, (iii) if Id
// generation requests are made at different milliseconds since epoch the length is guaranteed to be
// 9 symbols for 34 years with zero collisions, (iv) within the same millisecond unlimited number of
// further symbols can be added to guarantee no collisions, (v) 32 instead of 16 workers are
// supported.
//
// The algorithm uses less randomness than the original node.js implementation, which permits to
// extend the life span, reduce and guarantee the length. When encoding the worker and the
// millisecond value 2^5 (32) alphabet characters are used in 2 blocks for each symbol (original:
// 2^4 in 4 blocks). When encoding the count of requests within the same millisecond no randomness
// is used at all (i.e. one symbol represents 64 combinations). Therefore, multiple requests within
// the same millisecond will rarely lead to more than 2 extra symbols (request count must exceed
// 4096) and the value of 3 is difficult to imagine (260k requests within the same millisecond).
//
// The implemented type Abc exports the Encode method, which accepts 'digits' for n in 2^n of the
// alphabet characters to be used to encode data, thus permitting defining other algorithms with
// more or less randomness.
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

// Abc represents a shuffled alphabet used to generate the Ids and provides methods to
// encode data.
type Abc struct {
	alphabet []rune
}

// Shortid represents a short Id generator working with a given alphabet.
type Shortid struct {
	abc    Abc
	worker uint
	epoch  time.Time  // ids can be generated for 34 years since this date
	ms     uint       // ms since epoch for the last id
	count  uint       // request count within the same ms
	mx     sync.Mutex // locks access to ms and count
}

// DEFAULT_ABC is the default URL-friendly alphabet.
const DEFAULT_ABC = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

var shortid *Shortid

func init() {
	shortid = MustNew(0, DEFAULT_ABC, 1)
}

// GetDefault retrieves the default short Id generator initialised with the default alphabet,
// worker=0 and seed=1. The default can be overwritten using SetDefault.
func GetDefault() *Shortid {
	return shortid
}

// SetDefault overwrites the default generator.
func SetDefault(sid *Shortid) {
	shortid = sid
}

// Generate generates an Id using the default generator.
func Generate() (string, error) {
	return shortid.Generate()
}

// MustGenerate acts just like Generate, but panics instead of returning errors.
func MustGenerate() string {
	if id, err := Generate(); err == nil {
		return id
	} else {
		panic(err)
	}
}

// New constructs an instance of the short Id generator for the given worker number [0,31], alphabet
// (64 unique symbols) and seed value (to shuffle the alphabet). The worker number should be
// different for multiple or distributed processes generating Ids into the same data space. The
// seed, on contrary, should be identical.
func New(worker uint8, alphabet string, seed uint64) (*Shortid, error) {
	if abc, err := NewAbc(alphabet, seed); err == nil {
		sid := &Shortid{
			abc:    abc,
			worker: uint(worker),
			epoch:  time.Date(2016, time.January, 1, 0, 0, 0, 0, time.Local),
			ms:     0,
			count:  0,
		}
		log.Info("new %v", sid)
		return sid, nil
	} else {
		return nil, err
	}
}

// MustNew acts just like New, but panics instead of returning errors.
func MustNew(worker uint8, alphabet string, seed uint64) *Shortid {
	if sid, err := New(worker, alphabet, seed); err == nil {
		return sid
	} else {
		panic(err)
	}
}

// Generate generates a new short Id.
func (sid *Shortid) Generate() (string, error) {
	return sid.generate(time.Now(), sid.epoch)
}

// MustGenerate acts just like Generate, but panics instead of returning errors.
func (sid *Shortid) MustGenerate() string {
	if id, err := sid.Generate(); err == nil {
		return id
	} else {
		panic(err)
	}
}

func (sid *Shortid) generate(tm time.Time, epoch time.Time) (string, error) {
	ms, count := sid.msAndCountLocking(tm, epoch)
	idrunes := make([]rune, 9)
	if tmp, err := sid.abc.Encode(ms, 8, 5); err == nil {
		copy(idrunes, tmp) // first 8 symbols
	} else {
		return "", err
	}
	if tmp, err := sid.abc.Encode(sid.worker, 1, 5); err == nil {
		idrunes[8] = tmp[0]
	} else {
		return "", err
	}
	if count > 0 {
		if countrunes, err := sid.abc.Encode(count, 0, 6); err == nil {
			// only extend if really need it
			idrunes = append(idrunes, countrunes...)
		} else {
			return "", err
		}
	}
	return string(idrunes), nil
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

// String returns a string representation of the short Id generator.
func (sid *Shortid) String() string {
	return fmt.Sprintf("Shortid(worker=%v, epoch=%v, abc=%v)", sid.worker, sid.epoch, sid.abc)
}

// Abc returns the instance of alphabet used for representing the Ids.
func (sid *Shortid) Abc() Abc {
	return sid.abc
}

// Epoch returns the value of epoch used as the beginning of millisecond counting (normally
// 2016-01-01 00:00:00 local time)
func (sid *Shortid) Epoch() time.Time {
	return sid.epoch
}

// Worker returns the value of worker for this short Id generator.
func (sid *Shortid) Worker() uint {
	return sid.worker
}

// NewAbc constructs a new instance of shuffled alphabet to be used for Id representation.
func NewAbc(alphabet string, seed uint64) (Abc, error) {
	runes := []rune(alphabet)
	if len(runes) != len(DEFAULT_ABC) {
		return Abc{}, errors.New(fmt.Sprintf("alphabet must contain %v unique characters", len(DEFAULT_ABC)))
	}
	if nonUnique(runes) {
		return Abc{}, errors.New("alphabet must contain unique characters only")
	}
	abc := Abc{alphabet: nil}
	abc.shuffle(alphabet, seed)
	return abc, nil
}

// MustNewAbc acts just like NewAbc, but panics instead of returning errors.
func MustNewAbc(alphabet string, seed uint64) Abc {
	if res, err := NewAbc(alphabet, seed); err == nil {
		return res
	} else {
		panic(err)
	}
}

func nonUnique(runes []rune) bool {
	found := make(map[rune]struct{})
	for _, r := range runes {
		if _, seen := found[r]; !seen {
			found[r] = struct{}{}
		}
	}
	return len(found) < len(runes)
}

func (abc *Abc) shuffle(alphabet string, seed uint64) {
	source := []rune(alphabet)
	for len(source) > 1 {
		seed = (seed*9301 + 49297) % 233280
		i := int(seed * uint64(len(source)) / 233280)

		abc.alphabet = append(abc.alphabet, source[i])
		source = append(source[:i], source[i+1:]...)
	}
	abc.alphabet = append(abc.alphabet, source[0])
}

// Encode encodes a given value into a slice of runes of length nsymbols. In case nsymbols==0, the
// length of the result is automatically computed from data. Even if fewer symbols is required to
// encode the data than nsymbols, all positions are used encoding 0 where required to guarantee
// uniqueness in case further data is added to the sequence. The value of digits [4,6] represents
// represents n in 2^n, which defines how much randomness flows into the algorithm: 4 -- every value
// can be represented by 4 symbols in the alphabet (permitting at most 16 values), 5 -- every value
// can be represented by 2 symbols in the alphabet (permitting at most 32 values), 6 -- every value
// is represented by exactly 1 symbol with no randomness (permitting 64 values).
func (abc *Abc) Encode(val, nsymbols, digits uint) ([]rune, error) {
	if digits < 4 || 6 < digits {
		return nil, errors.New(fmt.Sprintf("allowed digits range [4,6], found %v", digits))
	}

	var computedSize uint = 1
	if val >= 1 {
		computedSize = uint(math.Log2(float64(val)))/digits + 1
	}
	if nsymbols == 0 {
		nsymbols = computedSize
	} else if nsymbols < computedSize {
		return nil, errors.New(fmt.Sprintf("cannot accommodate data, need %v digits, got %v", computedSize, nsymbols))
	}

	mask := 1<<digits - 1

	random := make([]int, int(nsymbols))
	// no random component if digits == 6
	if digits < 6 {
		copy(random, maskedRandomInts(len(random), 0x3f-mask))
	}

	res := make([]rune, int(nsymbols))
	for i, _ := range res {
		shift := digits * uint(i)
		index := (int(val>>shift) & mask) | random[i]
		res[i] = abc.alphabet[index]
	}
	return res, nil
}

// MustEncode acts just like Encode, but panics instead of returning errors.
func (abc *Abc) MustEncode(val, size, digits uint) []rune {
	if res, err := abc.Encode(val, size, digits); err == nil {
		return res
	} else {
		panic(err)
	}
}

func maskedRandomInts(size, mask int) []int {
	ints := make([]int, size)
	bytes := make([]byte, size)
	if _, err := randc.Read(bytes); err == nil {
		for i, b := range bytes {
			ints[i] = int(b) & mask
		}
	} else {
		for i, _ := range bytes {
			ints[i] = randm.Intn(0xff) & mask
		}
	}
	return ints
}

// String returns a string representation of the Abc instance.
func (abc Abc) String() string {
	return fmt.Sprintf("Abc{alphabet='%v')", abc.Alphabet())
}

// Alphabet returns the alphabet used as an immutable string.
func (abc Abc) Alphabet() string {
	return string(abc.alphabet)
}
