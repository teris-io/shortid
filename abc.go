package shortid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"time"
	"strings"
)

const ALPHABET = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

type Abc struct {
	alphabet  []rune
	startseed int64
	shuffled  []rune
	seed      int64
}

func NewAbc() (*Abc, error) {
	return NewAbcFor(nil, 1)
}

func MustNewAbc() *Abc {
	if abc, err := NewAbcFor(nil, 1); err == nil {
		return abc
	} else {
		panic(err)
	}
}

func NewAbcFor(rs []rune, seed int64) (*Abc, error) {
	if seed < 1 {
		return nil, errors.New("seed must be positive")
	}
	var alphabet []rune
	if rs == nil {
		alphabet = []rune(ALPHABET)
	} else {
		alphabet = UniqueRunes(rs)
		if len(alphabet) != len(ALPHABET) {
			return nil, errors.New(fmt.Sprintf("custom alphabet for shortid must be %v unique characters, found %v unique ones", len(ALPHABET), len(alphabet)))
		}
		SortRunes(alphabet)
	}
	return &Abc{alphabet: alphabet, startseed: seed, seed: seed}, nil
}

func MustNewAbcFor(rs []rune, seed int64) *Abc {
	if abc, err := NewAbcFor(rs, seed); err == nil {
		return abc
	} else {
		panic(err)
	}
}

func (abc *Abc) Shuffled() []rune {
	if len(abc.shuffled) == 0 {
		abc.shuffled = nil
		source := make([]rune, len(ALPHABET))
		copy(source, abc.alphabet)
		abc.next(len(source)) // copied from the original code, useless?
		for len(source) > 1 {
			i := abc.next(len(source))
			abc.shuffled = append(abc.shuffled, source[i])
			source = append(source[:i], source[i+1:]...)
		}
		abc.shuffled = append(abc.shuffled, source[0])
	}
	return abc.shuffled
}

func (abc *Abc) Reset() {
	abc.shuffled = nil
	abc.seed = abc.startseed
}

func (abc *Abc) Lookup(i int) (rune, error) {
	if i < 0 || len(ALPHABET) <= i {
		return 0, errors.New("index out of range")
	}
	return abc.Shuffled()[i], nil
}

func (abc *Abc) SetSeed(seed int64) {
	abc.startseed = seed
	abc.seed = seed
}

func (abc *Abc) Decode(id string) (uint8, uint8) {
	shuf := string(abc.Shuffled())
	rs := []rune(id)
	i1 :=strings.IndexRune(shuf, rs[0])
	i2 :=strings.IndexRune(shuf, rs[1])

	return uint8(i1) &0x0f, uint8(i2)&0x0f
}

func (abc *Abc) Encode(n uint) ([]rune, error) {
	buf := make([]byte, 1)
	var rs []rune

	maxi := int(math.Log2(float64(n)) -5)

	for i := 0; ; i++ {
		_, err := rand.Read(buf)
		if err != nil {
			return nil, err
		}
		// str = str + lookup( ( (number >> (4 * loopCounter)) & 0x0f ) | randomByte() );
		p1 := uint8(n>>uint(4*i))&0x0f
		p2 := uint8(buf[0])&0x30
		p3 := int(p1 | p2)
		p4 := uint(p3)
		if n == 5 {
			log.Debug("crypto random %v: %v | %v -> %v (%v)", n, p1, p2, p3, p4)
		}
		r, err := abc.Lookup(p3)
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
		if (maxi < i) {
			break
		}
	}
	return rs, nil
}

const worker = 0

var start time.Time
var version uint

func init() {
	start = time.Date(time.Now().Year()-1, time.November, 1, 1, 1, 1, 1, time.Local)
	version = uint(time.Now().Year() - 2016)
}

var lastSec uint = 0
var counter uint = 0

func (abc *Abc) Generate() (string, error) {

	sec := uint(time.Now().Sub(start).Seconds() *1e-3)
	log.Notice("sec %v", sec)

	if sec == lastSec {
		counter++;
	} else {
		counter = 0;
		lastSec = sec;
	}

	log.Debug("counter %v", counter)
	var res []rune

	var rs, err = abc.Encode(version)
	if err != nil {
		return "", err
	}
	res = append(res, rs...)

	rs, err = abc.Encode(worker)
	if err != nil {
		return "", err
	}
	res = append(res, rs...)

	if (counter > 0) {
		rs, err = abc.Encode(counter)
		if err != nil {
			return "", err
		}
		res = append(res, rs...)
	}
	rs, err = abc.Encode(sec)
	if err != nil {
		return "", err
	}
	res = append(res, rs...)


	return string(res), nil
}

// Based on The Central Randomizer 1.3 (C) 1997 by Paul Houle (houle@msc.cornell.edu)
func (abc *Abc) next(lessthen int) int {
	abc.seed = (abc.seed*9301 + 49297) % 233280
	return int(math.Floor(float64(abc.seed) / (233280.0) * float64(lessthen)))
}
