package icfp13

import (
	s "github.com/eadmund/sexprs"
	"math/rand"
	"strconv"
)

func GenVar(v Vars, n int) string {
	i := 0
	for k, v := range v {
		if i == n {
			return k
		}
		if v {
			i++
		}
	}
	panic("genvars should not reach this point")
}

func GenNewVar(v Vars) string {
	const prefix = "a__"
	for {
		n := rand.Int()
		name := prefix + strconv.Itoa(n)
		if !v[name] {
			return name
		}
	}
	panic("GenNewVar should never reach this point")
}

func GenSexp(e s.Sexp, v Vars) s.Sexp {
	muts := []Mutator{GenAtom, GenOp1s, GenOp2s, GenIf0, GenFold}
	l := len(muts)
	return muts[rand.Intn(l)](e, v)
}

func GenAtom(e s.Sexp, v Vars) s.Sexp {
	l := len(v)
	aname := string(e.(s.Atom).Value)
	for {
		i := Uint32n(l)
		n := GenVar(v, int(i))
		if n != aname {
			return MkAtom(n)
		}
	}
}

func GenOp1s(e s.Sexp, v Vars) s.Sexp {
	op1s := []string{"not", "shl1", "shr1", "shr4", "shr16"}
	op1l := len(op1s)
	body := GenAtom(e, v)
	op := op1s[rand.Intn(op1l)]
	return s.List{MkAtom(op), body}
}

func GenOp2s(e s.Sexp, v Vars) s.Sexp {
	op2s := []string{"and", "or", "xor", "plus"}
	op2l := len(op2s)
	left := GenAtom(e, v)
	right := GenAtom(e, v)
	op := op2s[rand.Intn(op2l)]
	return s.List{MkAtom(op), left, right}
}

func GenIf0(e s.Sexp, v Vars) s.Sexp {
	p := GenAtom(e, v)
	z := GenAtom(e, v)
	nz := GenAtom(e, v)
	return s.List{MkAtom("if0"), p, z, nz}
}

func GenFold(e s.Sexp, v Vars) s.Sexp {
	vec := GenAtom(e, v)
	start := GenAtom(e, v)
	arg1 := GenNewVar(v)
	v = v.Add(arg1)
	arg2 := GenNewVar(v)
	v = v.Add(arg2)
	body := GenAtom(e, v)
	return s.List{MkAtom("fold"), vec, start, s.List{MkAtom("lambda"), s.List{MkAtom(arg1), MkAtom(arg2)}, body}}
}
