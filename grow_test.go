package icfp13

import (
	s "github.com/eadmund/sexprs"
	"testing"
)

func TestFoo(t *testing.T) {
}

func TestCountAllMutationPoints(t *testing.T) {
	data := []struct {
		in  string
		out int
	}{
		{"(not x)", 1},
		{"(not (not (not 1)))", 1},
		{"(shl1 y)", 1},
		{"(shr1 z)", 1},
		{"(shr4 0)", 1},
		{"(shr16 u)", 1},
		{"(or x y)", 2},
		{"(and (and x (and y z)) (or a b))", 5},
		{"(xor (or 1 0) (and 1 1))", 4},
		{"(plus (not (xor 1 0)) (or (not 0) 1))", 4},
		{"(if0 a b c)", 3},
		{"(if0 (or a b) (not c) (plus (and d e) f))", 6},
		{"(fold v s (lambda (x y) e))", 3},
		{"(fold (not 1) (or x y) (lambda (a b) (plus a (and 0 1))))", 6},
	}

	for _, d := range data {
		e := Parse([]byte(d.in))
		if r := CountMutationPoints(e); r != d.out {
			t.Errorf("expected %v, got %v", d.out, r)
		}
	}
}

func TestMutateAt(t *testing.T) {
	data := []struct {
		in  string
		w   int
		out string
	}{
		{"(not x)", 1, "(not M)"},
		{"(not (or x y))", 2, "(not (or x M))"},
		{"(shl1 1)", 1, "(shl1 M)"},
		{"(shr1 1)", 1, "(shr1 M)"},
		{"(shr4 0)", 1, "(shr4 M)"},
		{"(shr16 b)", 1, "(shr16 M)"},
		{"(or a b)", 1, "(or M b)"},
		{"(or c d)", 2, "(or c M)"},
		{"(not (or a b))", 2, "(not (or a M))"},
		{"(or (or a b) (or b c))", 4, "(or (or a b) (or b M))"},
		{"(or (plus x 0) (xor b c))", 2, "(or (plus x M) (xor b c))"},
		{"(if0 a b c)", 1, "(if0 M b c)"},
		{"(if0 a b c)", 2, "(if0 a M c)"},
		{"(if0 a b c)", 3, "(if0 a b M)"},
		{"(fold 0 r (lambda (a b) (plus a b)))", 1, "(fold M r (lambda (a b) (plus a b)))"},
		{"(fold r 0 (lambda (a b) (plus a b)))", 2, "(fold r M (lambda (a b) (plus a b)))"},
		{"(fold x y (lambda (a b) (plus a b)))", 3, "(fold x y (lambda (a b) (plus M b)))"},
		{"(fold x y (lambda (a b) (plus a b)))", 4, "(fold x y (lambda (a b) (plus a M)))"},
	}

	f := func(e s.Sexp, vars []string) s.Sexp {
		return MkAtom("M")
	}

	for _, d := range data {
		e := Parse([]byte(d.in))
		ne, err := MutateAt(e, d.w, f)
		if err != nil {
			t.Errorf("failed with error: '%v'", err)
		} else if nes := ne.String(); d.out != nes {
			t.Errorf("Expected '%v', got '%v'", d.out, nes)
		}
	}
}

/*func TestGrowAtom(t *testing.T) {
  data := []struct {
    in s.Sexp
    env []string
    out s.Sexp
  } {
    { MkAtom("__1"), []string{"x"}, MkAtom("x") },
    { MkAtom("__0"), []string{"y"}, MkAtom("y") },
    { MkAtom("__0"), []string{}, MkAtom("__1") },
    { MkAtom("__1"), []string{}, MkAtom("__0") },
    { MkAtom("x"), []string{"x"}, MkAtom("x") },
    { MkAtom("x"), []string{"x","y"}, MkAtom("y") },
  }

  for _, d := range data {
    if newAtom := GrowAtom(d.in,d.env); !newAtom.Equal(d.out) {
      t.Errorf("expected %v, got '%v'", d.out, newAtom)
    }
  }
}*/
