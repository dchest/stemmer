Stemmer package for Go
======================

Includes `porter2` sub-package which implements English (Porter2) stemmer, as described by <http://snowball.tartarus.org/algorithms/english/stemmer.html>

Installation
-------------

    go get github.com/dchest/stemmer/porter2

This will install both the top-level `stemmer` and `stemmer/porter2` packages.

Example
-------

    import "github.com/dchest/stemmer/porter2"

    st := porter2.Stemmer
    st.Stem("delicious")   // => delici
    st.Stem("deliciously") // => delici

Tests
-----

porter2:

Included `test_output.txt` and `test_voc.txt` are from [the original implementation](http://snowball.tartarus.org/algorithms/english/stemmer.html), used only when running tests with `go test`.
