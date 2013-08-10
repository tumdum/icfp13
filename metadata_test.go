package icfp13

import (
	s "github.com/eadmund/sexprs"
	"testing"
)

func TestAtomSize(t *testing.T) {
	data := []s.Atom{MkAtom("x"), MkAtom("__1"), MkAtom("__0")}
	for _, d := range data {
		if size := Size(d); size != 1 {
			t.Errorf("expected size 1, got %d", size)
		}
	}
}

func TestListSize(t *testing.T) {
	data := []struct {
		in  string
		out int
	}{
		{"(if0 0 0 0)", 4},
		{"(if0 (if0 0 0 0) 0 0)", 7},
		{"(fold e0 e1 (lambda (x y) e2))", 5},
		{"(fold (if0 0 0 0) (if0 0 0 0) (lambda (x y) (if0 0 0 0)))", 14},
		{"(lambda (x) (if0 0 0 0))", 5},
		{"(or 0 0)", 3},
		{"(shr16 1)", 2},
	}

	for _, d := range data {
		e := Parse([]byte(d.in))
		if r := Size(e); r != d.out {
			t.Errorf("expected %v got %v", d.out, r)
		}
	}
}
