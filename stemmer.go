// Copyright (c) 2011 Dmitry Chestnykh
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package stemmer declares Stemmer interface.
package stemmer

type Stemmer interface{
	Stem(s string) string
}