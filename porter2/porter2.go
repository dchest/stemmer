// Copyright (c) 2011 Dmitry Chestnykh
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package porter2 implements English (Porter2) stemmer, as described by
// http://snowball.tartarus.org/algorithms/english/stemmer.html
package porter2

import (
	"strings"
	"github.com/dchest/stemmer"
)

// Stemmer is a global, shared instance of Porter2 English stemmer.
var Stemmer stemmer.Stemmer

type englishStemmer bool

func init() {
	Stemmer = englishStemmer(true)
}

func suffixPos(s, suf []int) int {
	if len(s) < len(suf) {
		return -1
	}
	j := len(s) - 1
	for i := len(suf) - 1; i >= 0; i-- {
		if suf[i] != s[j] {
			return -1
		}
		j--
	}
	return len(s) - len(suf)
}

func removeSuffix(s, suf []int) []int {
	i := suffixPos(s, suf)
	if i != -1 {
		return s[:i]
	}
	return s
}

func isVowel(rune int) bool {
	switch rune {
	case 'a', 'e', 'i', 'o', 'u', 'y':
		return true
	}
	return false
}

func hasVowelBeforePos(s []int, pos int) bool {
	for i := pos; i >= 0; i-- {
		if isVowel(s[i]) {
			return true
		}
	}
	return false
}

var rExceptions = []string{
	"gener",
	"commun",
	"arsen",
}

func calcR(s []int) int {
	for i := 0; i < len(s)-1; i++ {
		if isVowel(s[i]) && !isVowel(s[i+1]) {
			return i + 2
		}
	}
	return len(s)
}

func getR1R2(s []int) (r1, r2 int) {
	for _, v := range rExceptions {
		if strings.HasPrefix(string(s), v) {
			r1 = len(v)
			r2 = r1 + calcR(s[r1:])
			return
		}
	}
	r1 = calcR(s)
	r2 = r1 + calcR(s[r1:])
	return
}

func endsWithDouble(s []int) bool {
	if len(s) < 2 {
		return false
	}
	last := s[len(s)-1]
	switch last {
	case 'b', 'd', 'f', 'g', 'm', 'n', 'p', 'r', 't':
		if s[len(s)-2] == last {
			return true
		}
	}
	return false
}

func isShortWord(s []int) bool {
	if r1, _ := getR1R2(s); r1 != len(s) {
		return false
	}
	i := len(s)
	if i == 2 && isVowel(s[0]) && !isVowel(s[1]) {
		return true
	}
	if i < 3 {
		return false
	}
	// ends with short sillable?
	// N + v + N
	last := s[i-1]
	if !isVowel(s[i-3]) && isVowel(s[i-2]) && !isVowel(last) &&
		last != 'w' && last != 'x' && last != 'Y' {
		return true
	}
	return false
}

var step1bWords = [][]int{
	[]int("ingly"),
	[]int("edly"),
	[]int("ing"),
	[]int("ed"),
}

var step2Words = [][]int{
	[]int("fulness"), // ful
	[]int("ousness"), // ous
	[]int("iveness"), // ive
	[]int("ational"), // ate
	[]int("ization"), // ize
	[]int("tional"),  // tion
	[]int("biliti"),  // ble
	[]int("lessli"),  // less
	[]int("fulli"),   // ful
	[]int("ousli"),   // ous
	[]int("iviti"),   // ive
	[]int("alism"),   // al
	[]int("ation"),   // ate
	[]int("entli"),   // ent
	[]int("aliti"),   // al
	[]int("enci"),    // ence
	[]int("anci"),    // ance
	[]int("abli"),    // able
	[]int("izer"),    // ize
	[]int("ator"),    // ate
	[]int("alli"),    // al
	[]int("bli"),     // ble
	//"ogi",   // replace with og if preceded by l -- handled later in code
	//"li"     // delete if preceded by a valid li-ending  -- handled later code
}

var step2Reps = [][]int{
	[]int("ful"),
	[]int("ous"),
	[]int("ive"),
	[]int("ate"),
	[]int("ize"),
	[]int("tion"),
	[]int("ble"),
	[]int("less"),
	[]int("ful"),
	[]int("ous"),
	[]int("ive"),
	[]int("al"),
	[]int("ate"),
	[]int("ent"),
	[]int("al"),
	[]int("ence"),
	[]int("ance"),
	[]int("able"),
	[]int("ize"),
	[]int("ate"),
	[]int("al"),
	[]int("ble"),
	//"og"  -- handled later in code
	// ""   -- handled later in code
}

var step3Words = [][]int{
	[]int("ational"), // ate
	[]int("tional"),  // tion
	[]int("alize"),   // al
	[]int("icate"),   // ic
	[]int("iciti"),   // ic
	[]int("ical"),    // ic
	[]int("ful"),     // (delete)
	[]int("ness"),    // (delete)
	//ative -- handled later in code
}

var step3Reps = [][]int{
	[]int("ate"),
	[]int("tion"),
	[]int("al"),
	[]int("ic"),
	[]int("ic"),
	[]int("ic"),
	[]int{},
	[]int{},
	[]int{},
}

var step4Words = [][]int{
	[]int("ement"),
	[]int("able"),
	[]int("ible"),
	[]int("ance"),
	[]int("ence"),
	[]int("ment"),
	[]int("ant"),
	[]int("ent"),
	[]int("ism"),
	[]int("ate"),
	[]int("iti"),
	[]int("ous"),
	[]int("ive"),
	[]int("ize"),
	[]int("al"),
	[]int("er"),
	[]int("ic"),
	// "ion" -- delete if preceded by s or t
}

var exceptions1 = map[string]string{
	// special changes
	"skis":  "ski",
	"skies": "sky",
	"dying": "die",
	"lying": "lie",
	"tying": "tie",

	// special -LY cases
	"idly":   "idl",
	"gently": "gentl",
	"ugly":   "ugli",
	"early":  "earli",
	"only":   "onli",
	"singly": "singl",
	//invariant forms
	"sky":  "sky",
	"news": "news",
	"howe": "howe",
	// not plural forms
	"atlas":  "atlas",
	"cosmos": "cosmos",
	"bias":   "bias",
	"andes":  "andes",
}

var exceptions2 = map[string]bool{
	"inning":  true,
	"outing":  true,
	"canning": true,
	"herring": true,
	"earring": true,
	"proceed": true,
	"exceed":  true,
	"succeed": true,
}

// Stem returns a stemmed string word
func (stm englishStemmer) Stem(word string) string {
	word = strings.ToLower(word)
	// Is it exception?
	if rep, ex := exceptions1[word]; ex {
		return rep
	}
	s := []int(word)
	if len(s) <= 2 {
		return word
	}
	if s[0] == '\'' {
		s = s[1:]
	}
	if s[0] == 'y' {
		s[0] = 'Y'
	}
	for i := 1; i < len(s); i++ {
		if isVowel(s[i-1]) && s[i] == 'y' {
			s[i] = 'Y'
		}
	}
	r1, r2 := getR1R2(s)

	// Step 0
	s = removeSuffix(s, []int("'s'"))
	s = removeSuffix(s, []int("'s"))
	s = removeSuffix(s, []int("'"))

	// Step 1a
	if i := suffixPos(s, []int("sses")); i != -1 {
		// sses
		// replace by ss
		s = append(s[:i], []int("ss")...)
		goto step1b
	}
	{
		i := suffixPos(s, []int("ied"))
		if i == -1 {
			i = suffixPos(s, []int("ies"))
		}
		if i != -1 {
			// ied+   ies*
			// replace by i if preceded by more than one letter,
			// otherwise by ie (so ties -> tie, cries -> cri)
			s = s[:i]
			if len(s) > 1 {
				s = append(s, int('i'))
			} else {
				s = append(s, []int("ie")...)
			}
			goto step1b
		}
	}
	if suffixPos(s, []int("us")) != -1 || suffixPos(s, []int("ss")) != -1 {
		// do nothing
		goto step1b
	}

	if i := suffixPos(s, []int("s")); i != -1 {
		if len(s) >= 3 && hasVowelBeforePos(s, len(s)-3) {
			s = s[:i]
		}
		goto step1b
	}

step1b:
	if _, ex := exceptions2[string(s)]; ex {
		return string(s)
	}
	// Step 1b
	for _, suf := range [][]int{[]int("eed"), []int("eedly")} {
		if i := suffixPos(s, suf); i != -1 {
			if i >= r1 {
				s = append(s[:i], []int("ee")...)
			}
			goto step1c
		}
	}

	for _, suf := range step1bWords {
		if suffixPos(s, suf) != -1 {
			if len(s) > len(suf) && hasVowelBeforePos(s, len(s)-len(suf)-1) {
				s = s[:len(s)-len(suf)]
			} else {
				goto step1c
			}
			if suffixPos(s, []int("at")) != -1 || suffixPos(s, []int("bl")) != -1 ||
				suffixPos(s, []int("iz")) != -1 {
				s = append(s, int('e'))
				goto step1c
			}
			if endsWithDouble(s) {
				s = s[:len(s)-1]
				goto step1c
			}
			if isShortWord(s) {
				s = append(s, int('e'))
			}
			goto step1c
		}
	}
step1c:
	// replace suffix y or Y by i if preceded by a non-vowel which is
	// not the first letter of the word (so cry -> cri, by -> by, say -> say)
	if len(s) > 2 {
		switch s[len(s)-1] {
		case 'y', 'Y':
			if !isVowel(s[len(s)-2]) {
				s[len(s)-1] = 'i'
			}
		}
	}
	//step2:
	r1, r2 = getR1R2(s)
	// Search for the longest among the following suffixes, and,
	// if found and in R1, perform the action indicated
	for j, suf := range step2Words {
		if i := suffixPos(s, suf); i != -1 {
			if i >= r1 {
				s = append(s[:i], step2Reps[j]...)
			}
			goto step3
		}
	}
	if i := suffixPos(s, []int("ogi")); i != -1 && i >= r1 {
		if s[i-1] == 'l' {
			s = append(s[:i], []int("og")...)
		}
		goto step3
	}
	if i := suffixPos(s, []int("li")); i != -1 && i >= r1 {
		// valid li-ending: c   d   e   g   h   k   m   n   r   t
		switch s[i-1] {
		case 'c', 'd', 'e', 'g', 'h', 'k', 'm', 'n', 'r', 't':
			s = s[:i]
		}
	}
step3:
	r1, r2 = getR1R2(s)
	for j, suf := range step3Words {
		if i := suffixPos(s, suf); i != -1 {
			if i >= r1 {
				s = append(s[:i], step3Reps[j]...)
			}
			goto step4
		}
	}
	if i := suffixPos(s, []int("ative")); i != -1 && i >= r2 {
		s = s[:i]
		goto step4
	}

step4:
	r1, r2 = getR1R2(s)
	for _, suf := range step4Words {
		if i := suffixPos(s, suf); i != -1 {
			if i >= r2 {
				s = s[:i]
			}
			goto step5
		}
	}
	if i := suffixPos(s, []int("ion")); i != -1 && i >= r2 {
		switch s[i-1] {
		case 's', 't':
			s = s[:i]
		}
	}

step5:
	r1, r2 = getR1R2(s)
	i := len(s) - 1
	if i > 0 && s[i] == 'e' {
		if i >= r2 {
			s = s[:i]
			goto final
		}
		if i >= r1 {
			// if not preceded by a short syllable
			if i < 3 {
				goto final
			}
			// N + v + N
			last := s[i-1]
			if !isVowel(s[i-3]) && isVowel(s[i-2]) && !isVowel(last) &&
				last != 'w' && last != 'x' && last != 'Y' {
				goto final
			}
			s = s[:i]
		}
		goto final
	}
	if i > 1 && i >= r2 && s[i] == 'l' && s[i-1] == 'l' {
		s = s[:i]
	}

final:
	for i, v := range s {
		if v == 'Y' {
			s[i] = 'y'
		}
	}
	return string(s)
}
