// Copyright 2016 Ventu.io. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

package shortid_test

import (
	"github.com/ventu-io/go-shortid"
	"math/rand"
	"testing"
)

func TestAbc_onSchuffle_seedDependentAndReproducible(t *testing.T) {
	abc := shortid.MustNewAbc()
	res := string(abc.Shuffled())
	if res != "ylZM7VHLvOFcohp01x-fXNr8P_tqin6RkgWGm4SIDdK5s2TAJebzQEBUwuY9j3aC" {
		t.Error("incorrect shuffled abc for seed 1")
	}

	abc = shortid.MustNewAbcFor(nil, 1234)
	res = string(abc.Shuffled())
	if res != "ef4w9iMboqLOQdWu3hKI72A0VZpCtzDlXk5_a6cFSNYGnH-gmsP1UBxvTRJjE8ry" {
		t.Error("incorrect shuffled abc for seed 1234")
	}
}

func TestAbc_onShuffle_runeSliceOfSameLengthAsABC_andUnique(t *testing.T) {
	seed := rand.Int63n(1024)
	abc := shortid.MustNewAbcFor(nil, seed)
	if len(shortid.UniqueRunes(abc.Shuffled())) != len(shortid.ALPHABET) {
		t.Error("incorrect length or not unique")
	}
}

func TestAbc_onConstruction_worksWithCustomAlphabets(t *testing.T) {

	// wrong length
	// non-unique

	// funky
	//alphabet.characters("①②③④⑤⑥⑦⑧⑨⑩⑪⑫ⒶⒷⒸⒹⒺⒻⒼⒽⒾⒿⓀⓁⓂⓃⓄⓅⓆⓇⓈⓉⓊⓋⓌⓍⓎⓏⓐⓑⓒⓓⓔⓕⓖⓗⓘⓙⓚⓛⓜⓝⓞⓟⓠⓡⓢⓣⓤⓥⓦⓧⓨⓩ");
	//expect(alphabet.shuffled()).to.equal("ⓌⒿⓧⓚ⑧ⓣⓕⓙⓉⓜⓓⒶⓂⒻⓃ①②ⓋⓩⒹⓥⓛⓅ⑨ⓝⓨⓇⓄⒼⓁ⑦ⓟⒾⒺⓤⓔⓀ⑤ⓠⓖⓑⒷⓘ⑥Ⓠ③ⓡⓎⓗⒸ⑫ⓍⓞⓒⓏⓢⓊⓈⓦ⑩Ⓗ④⑪ⓐ");
}

func TestAbc_Encode(t *testing.T) {
	abc := shortid.MustNewAbcFor(nil, 1234)
	for i := 0; i < 32; i++ {
		id, err := abc.Generate()
		if err != nil {
			t.Error(err)
		}
		v, w := abc.Decode(id)
		log.Notice("encoded %v as %v (%v/%v)", i, id, v, w)
	}
}
