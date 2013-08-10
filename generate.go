package icfp13

import (
	s "github.com/eadmund/sexprs"
	"math/rand"
	"strconv"
)

const NewGenerationSize = 300

func GenNewRandomUniqGeneration(start s.Sexp) []s.Sexp {
	return genNewRandomUniqGeneration(start, GenSexp, NewGenerationSize)
}

func GenNewRandomUniqGenerationUsing(start s.Sexp, ops []string, size int) []s.Sexp {
	return genNewRandomUniqGeneration(start, MetaMutator2(ops), size)
}

func genNewRandomUniqGeneration(start s.Sexp, m Mutator, maxSize int) []s.Sexp {
	muts := CountMutationPoints(start)
	ret := make([]s.Sexp, 0, NewGenerationSize)
	seen := make(map[string]bool)

	fails := 0
	for (len(ret) < NewGenerationSize) && (fails < NewGenerationSize/5) {
		mpoint := rand.Intn(muts) + 1
		mutation, err := MutateAt(start, mpoint, m)
		if err != nil {
			fails++
			continue
		}
		mstr := mutation.String()
		if seen[mstr] {
			fails++
			continue
		} else {
			ret = append(ret, mutation)
			seen[mstr] = true
		}
	}
	return ret
}

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
		n := len(v)
		name := prefix + strconv.Itoa(n)
		if !v[name] {
			return name
		}
	}
	panic("GenNewVar should never reach this point")
}

func RandMut(ops []Mutator) Mutator {
	return ops[rand.Intn(len(ops))]
}

func MetaMutator2(ops []string) Mutator {
	m := make(map[string]Mutator)
	m["not"] = MetaOp1Named("not")
	m["shl1"] = MetaOp1Named("shl1")
	m["shr1"] = MetaOp1Named("shr1")
	m["shr4"] = MetaOp1Named("shr4")
	m["shr16"] = MetaOp1Named("shr16")
	m["and"] = MetaOp2Named("and")
	m["or"] = MetaOp2Named("or")
	m["xor"] = MetaOp2Named("xor")
	m["plus"] = MetaOp2Named("plus")
	m["if0"] = GenIf0
	m["fold"] = GenFold
	m["tfold"] = GenFold

	selected := make([]Mutator, 0)
	for _, op := range ops {
		if v, ok := m[op]; ok {
			selected = append(selected, v)
		}
	}
	return func(e s.Sexp, v Vars) s.Sexp {
		return RandMut(selected)(e, v)
	}
}

func MetaMutator(ops []string) Mutator {
	muts := []Mutator{GenAtom, GenOp1s, GenOp2s, GenIf0}
	if len(ops) > 0 {
		muts = append(muts, GenFold)
	}
	return func(e s.Sexp, v Vars) s.Sexp {
		return RandMut(muts)(e, v)
	}
}

func GenSexp(e s.Sexp, v Vars) s.Sexp {
	muts := []Mutator{GenAtom, GenOp1s, GenOp2s, GenIf0, GenFold}
	return RandMut(muts)(e, v)
}

func GenAtom(e s.Sexp, v Vars) s.Sexp {
	l := len(v)
	for {
		i := Uint32n(l)
		n := GenVar(v, int(i))
		return MkAtom(n)
	}
}

func MetaOp1Named(op string) Mutator {
	return func(e s.Sexp, v Vars) s.Sexp {
		body := GenAtom(e, v)
		return s.List{MkAtom(op), body}
	}
}

func GenOp1s(e s.Sexp, v Vars) s.Sexp {
	op1s := []string{"not", "shl1", "shr1", "shr4", "shr16"}
	op1l := len(op1s)
	body := GenAtom(e, v)
	op := op1s[rand.Intn(op1l)]
	return s.List{MkAtom(op), body}
}

func MetaOp2Named(op string) Mutator {
	return func(e s.Sexp, v Vars) s.Sexp {
		l := GenAtom(e, v)
		r := GenAtom(e, v)
		return s.List{MkAtom(op), l, r}
	}
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
