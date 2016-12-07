// Copyright 2016 The Stemmer Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dutch implements Dutch stemmer, as described in
// http://snowball.tartarus.org/algorithms/dutch/stemmer.html
package dutch

import (
	"strings"

	"github.com/dchest/stemmer"
)

// Stemmer is a global, shared instance of Dutch stemmer.
var Stemmer stemmer.Stemmer = dutchStemmer(true)

type dutchStemmer bool

func suffixPos(s, suf []rune) int {
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

func firstSuffixPos(s []rune, suffixes ...[]rune) int {
	for _, suf := range suffixes {
		if i := suffixPos(s, suf); i >= 0 {
			return i
		}
	}
	return -1
}

func isVowel(x rune) bool {
	switch x {
	case 'a', 'e', 'i', 'o', 'u', 'y', 'è':
		return true
	default:
		return false
	}
}

func isVowelAt(s []rune, index int) bool {
	if index < 0 || index >= len(s) {
		return false
	}
	return isVowel(s[index])
}

func calcR(s []rune) int {
	for i := 0; i < len(s)-1; i++ {
		if isVowelAt(s, i) && !isVowelAt(s, i+1) {
			return i + 2
		}
	}
	return len(s)
}

func adjustR1(s []rune, r1 int) int {
	if r1 >= 3 {
		return r1
	}
	if len(s) < 4 {
		return len(s)
	}
	return 3
}

func getR1R2(s []rune) (r1, r2 int) {
	r1 = calcR(s)
	r2 = r1 + calcR(s[r1:])
	r1 = adjustR1(s, r1)
	return
}

func hasValidSEnding(s []rune) bool {
	last := s[len(s)-1]
	if last == 'j' {
		return false
	}
	return !isVowel(last)
}

func hasValidEnEnding(s []rune) bool {
	last := s[len(s)-1]
	if isVowel(last) {
		return false
	}
	return suffixPos(s, []rune("gem")) < 0
}

// Delete suffix e if in R1 and preceded by a non-vowel
func deleteESuffixPrecededByNonVowel(s []rune, r1 int) ([]rune, bool) {
	if i := len(s) - 1; i >= r1 && s[i] == 'e' && !isVowelAt(s, i-1) {
		s = s[:i]
		return s, true
	}
	return s, false
}

func undouble(s []rune) []rune {
	if len(s) < 2 {
		return s
	}
	last := s[len(s)-1]
	switch last {
	case 'k', 'd', 't':
		prev := s[len(s)-2]
		if prev == last {
			return s[:len(s)-1]
		}
	}
	return s
}

// Stem returns a stemmed string word.
func (stm dutchStemmer) Stem(word string) string {
	word = strings.ToLower(word)
	s := []rune(word)
	for i, c := range s {
		switch c {
		case 'ä':
			s[i] = 'a'
		case 'ë':
			s[i] = 'e'
		case 'ï':
			s[i] = 'i'
		case 'ö':
			s[i] = 'o'
		case 'ü':
			s[i] = 'u'
		case 'á':
			s[i] = 'a'
		case 'é':
			s[i] = 'e'
		case 'í':
			s[i] = 'i'
		case 'ó':
			s[i] = 'o'
		case 'ú':
			s[i] = 'u'
		}
	}

	// Put initial y into uppercase
	if len(s) > 0 && s[0] == 'y' {
		s[0] = 'Y'
	}
	// Put y after a vowel into upper case
	for i, max := 1, len(s); i < max; i++ {
		if s[i] == 'y' && isVowelAt(s, i-1) {
			s[i] = 'Y'
		}
	}
	// Put i between vowels into upper case
	for i, max := 1, len(s)-1; i < max; i++ { // interested only in runes between vowels, so we run from 1 to len-1
		if s[i] == 'i' && isVowelAt(s, i-1) && isVowelAt(s, i+1) {
			s[i] = 'I'
		}
	}

	r1, r2 := getR1R2(s)

	// step 1 group a
	if i := suffixPos(s, []rune("heden")); i >= 0 {
		if i >= r1 {
			// replace with heid if in R1
			s = append(s[:i], []rune("heid")...)
		}
		goto step2
	}

	// step 1 group b
	if i := firstSuffixPos(s, []rune("ene"), []rune("en")); i >= 0 {
		// delete if in R1 and preceded by a valid en-ending, and then undouble the ending
		if i >= r1 && hasValidEnEnding(s[:i]) {
			s = s[:i]
			s = undouble(s)
		}
		goto step2
	}

	// step 1 group c
	if i := firstSuffixPos(s, []rune("se"), []rune("s")); i >= 0 {
		// delete if in R1 and preceded by a valid s-ending
		if i >= r1 && hasValidSEnding(s[:i]) {
			s = s[:i]
		}
		goto step2
	}

step2:
	// step 2 group a
	// Delete suffix e if in R1 and preceded by a non-vowel, and then undouble the ending
	s, step2RemovedE := deleteESuffixPrecededByNonVowel(s, r1)
	if step2RemovedE {
		s = undouble(s)
	}

	// step 3 heid
	if i := suffixPos(s, []rune("heid")); i >= r2 && s[i-1] != 'c' {
		// delete heid if in R2 and not preceded by c, and treat a preceding en as in step 1(b)
		prefix := s[:i]
		suffix := s[i+4:]

		if i := suffixPos(prefix, []rune("en")); i >= r1 && hasValidEnEnding(prefix[:i]) {
			// delete if in R1 and preceded by a valid en-ending, and then undouble the ending
			prefix = prefix[:i]
			prefix = undouble(prefix)
		}

		s = append(prefix, suffix...)
	}

	// step 3b: end ing
	if i := firstSuffixPos(s, []rune("end"), []rune("ing")); i >= 0 {
		if i >= r2 {
			// delete if in R2
			s = s[:i]
			// if preceded by ig, delete if in R2 and not preceded by e, otherwise undouble the ending
			if i = suffixPos(s, []rune("ig")); i >= r2 && suffixPos(s[:i], []rune("e")) == -1 {
				s = s[:i]
			} else {
				s = undouble(s)
			}
		}
		goto step4
	}

	// step 3 ig
	if i := suffixPos(s, []rune("ig")); i >= 0 {
		if i >= r2 && suffixPos(s[:i], []rune("e")) == -1 {
			s = s[:i]
		}
		goto step4
	}

	// step 3 lijk
	if i := suffixPos(s, []rune("lijk")); i >= 0 {
		if i >= r2 {
			s = s[:i]
			s, _ = deleteESuffixPrecededByNonVowel(s, r1)
			s = undouble(s)
		}
		goto step4
	}

	// step 3 baar
	if i := suffixPos(s, []rune("baar")); i >= 0 {
		if i >= r2 {
			s = s[:i]
		}
		goto step4
	}

	// step 3 bar
	if i := suffixPos(s, []rune("bar")); i >= 0 {
		if step2RemovedE && i >= r2 {
			s = s[:i]
		}
		goto step4
	}

step4: // undouble vowel
	// If the words ends CVD, where C is a non-vowel, D is a non-vowel other than I, and V is double a, e, o or u,
	// remove one of the vowels from V (for example, maan -> man, brood -> brod).
	l := len(s)
	if l >= 4 {
		c, v1, v2, d := s[l-4], s[l-3], s[l-2], s[l-1]
		if v1 == v2 {
			switch v1 {
			case 'a', 'e', 'o', 'u':
				if !isVowel(c) && !isVowel(d) && d != 'I' {
					s = append(s[:l-2], s[l-1:]...)
				}
			}
		}
	}

	// finally
	for i, max := 0, len(s); i < max; i++ {
		switch s[i] {
		case 'Y':
			s[i] = 'y'
		case 'I':
			s[i] = 'i'
		}
	}

	return string(s)
}
