// Copyright 2012 The Stemmer Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package german

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func TestStem(t *testing.T) {
	voc, err := os.Open("test_voc.txt")
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	defer voc.Close()
	output, err := os.Open("test_output.txt")
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	defer output.Close()
	bvoc := bufio.NewReader(voc)
	bout := bufio.NewReader(output)
	for {
		vocline, err := bvoc.ReadString('\n')
		if err != nil {
			switch err {
			case io.EOF:
				return
			default:
				t.Errorf("%s", err)
				return
			}
		}
		outline, err := bout.ReadString('\n')
		if err != nil {
			switch err {
			case io.EOF:
				return
			default:
				t.Errorf("%s", err)
				return
			}
		}
		vocline = vocline[:len(vocline)-1]
		outline = outline[:len(outline)-1]
		st := Stemmer.Stem(vocline)
		if st != outline {
			t.Errorf("\"%s\" expected %q got %q", vocline, outline, st)
		}
	}
}
