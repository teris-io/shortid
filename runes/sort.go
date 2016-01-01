// Copyright (c) deft_code 2014
// (http://stackoverflow.com/users/28817/deft-code)
//
// Adapted from: http://stackoverflow.com/questions/22688651

package runes

import syssort "sort"

type sortrunes []rune

func (s sortrunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortrunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortrunes) Len() int {
	return len(s)
}

func Sort(rs []rune) {
	syssort.Sort(sortrunes(rs))
}

func Unique(rs []rune) []rune {
	var res []rune
	found := make(map[rune]struct{})
	for _, r := range rs {
		if _, seen := found[r]; !seen {
			found[r] = struct{}{}
			res = append(res, r)
		}
	}
	return res
}
