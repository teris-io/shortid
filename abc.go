package shortid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/ventu-io/go-shortid/runes"
	"math"
	"time"
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
		alphabet = runes.Unique(rs)
		if len(alphabet) != len(ALPHABET) {
			return nil, errors.New(fmt.Sprintf("custom alphabet for shortid must be %v unique characters, found %v unique ones", len(ALPHABET), len(alphabet)))
		}
		runes.Sort(alphabet)
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

func (abc *Abc) Decode(id string) (int8, int8) {
	rs := []rune(id)
	var vers int8 = -1
	var worker int8 = -1
	for i, r := range abc.Shuffled() {
		if vers >= 0 && worker >= 0 {
			break
		}
		if rs[0] == r {
			vers = int8(i) & 0x0f
		}
		if rs[1] == r {
			worker = int8(i) & 0x0f
		}
	}
	return vers, worker
}

func (abc *Abc) Encode(n uint) (string, error) {
	buf := make([]byte, 1)
	var rs []rune

	maxi := int(math.Log2(float64(n)) -5)

	for i := 0; ; i++ {
		_, err := rand.Read(buf)
		if err != nil {
			return "", err
		}
		r, err := abc.Lookup(int((n>>uint(4*i))&0x0f) | int(buf[0]&0x30))
		if err != nil {
			return "", err
		}
		rs = append(rs, r)
		if (maxi < i) {
			break
		}
	}
	return string(rs), nil
}

const version  = 5
const worker = 0

var start time.Time

func init() {
	start = time.Date(time.Now().Year()-1, time.November, 1, 1, 1, 1, 1, time.Local)
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

	var str, err = abc.Encode(version)
	if err != nil {
		return "", err
	}
	res := str

	str, err = abc.Encode(worker)
	if err != nil {
		return "", err
	}
	res += str

	if (counter > 0) {
		str, err = abc.Encode(counter)
		if err != nil {
			return "", err
		}
		res += str
	}
	str, err = abc.Encode(sec)
	if err != nil {
		return "", err
	}
	res += str


	return str, nil
}

// Based on The Central Randomizer 1.3 (C) 1997 by Paul Houle (houle@msc.cornell.edu)
func (abc *Abc) next(lessthen int) int {
	abc.seed = (abc.seed*9301 + 49297) % 233280
	return int(math.Floor(float64(abc.seed) / (233280.0) * float64(lessthen)))
}
