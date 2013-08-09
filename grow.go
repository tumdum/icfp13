package icfp13

import (
	"math/rand"
	//"fmt"
	"errors"
	s "github.com/eadmund/sexprs"
)

const StartingSexp = "(lambda (x) x)"

func MkAtom(name string) s.Atom {
	return s.Atom{[]byte{}, []byte(name)}
}

func Uint32n(n int) uint32 {
	return rand.Uint32() % uint32(n)
}

func GenVar(vars []string) s.Sexp {
	i := Uint32n(len(vars))
	return MkAtom(vars[i])
}

func CountMutationPoints(e s.Sexp) int {
	switch e := e.(type) {
	case s.List:
		return countMutationPointsInList(e)
	case s.Atom:
		return 1
	default:
		panic("?")
	}
}

func countMutationPointsInList(l s.List) int {
	head := string(l[0].(s.Atom).Value)
	switch head {
	case "not", "shl1", "shr1", "shr4", "shr16":
		return CountMutationPoints(l[1])
	case "or", "and", "xor", "plus":
		return CountMutationPoints(l[1]) + CountMutationPoints(l[2])
	case "if0":
		return CountMutationPoints(l[1]) + CountMutationPoints(l[2]) + CountMutationPoints(l[3])
	case "lambda":
		return CountMutationPoints(l[2])
	case "fold":
		return countMutationPointsInFold(l[1], l[2], l[3])
	default:
		return 0
	}
}

func countMutationPointsInFold(vec, start, lambda s.Sexp) int {
	l := lambda.(s.List)
	lbody := l[2]
	return CountMutationPoints(vec) + CountMutationPoints(start) + CountMutationPoints(lbody)
}

type Mutator func(s.Sexp, []string) s.Sexp

func MutateAt(e s.Sexp, where int, m Mutator) (s.Sexp, error) {
	// fmt.Println("---------")
	me, n := mutateAt(e, where, m, []string{})
	if n != 0 {
		return nil, errors.New("failed to mutate")
	}
	return me, nil
}

func mutateAt(e s.Sexp, where int, m Mutator, vars []string) (s.Sexp, int) {
	switch e := e.(type) {
	case s.List:
		return mutateList(e, where, m, vars)
	case s.Atom:
		return mutateAtom(e, where, m, vars)
	default:
		panic("mutate unknown type")
	}
}

func mutateAtom(a s.Atom, where int, m Mutator, vars []string) (s.Sexp, int) {
	// fmt.Println("mutateAtom", a, where)
	if where == 1 {
		return m(a, vars), 0
	}
	if where > 0 {
		where--
	}
	return MkAtom(string(a.Value)), where
}

func mutateList(l s.List, where int, m Mutator, vars []string) (s.Sexp, int) {
	head := string(l[0].(s.Atom).Value)
	switch head {
	case "not", "shl1", "shr1", "shr4", "shr16":
		return mutateOp1(l, where, m, vars)
	case "or", "and", "xor", "plus":
		return mutateOp2(l, where, m, vars)
	case "if0":
		return mutateIf0(l, where, m, vars)
	case "lambda":
		return mutateLambda(l, where, m, vars)
	case "fold":
		return mutateFold(l, where, m, vars)
	default:
		panic("mutate list with unknown head: " + head)
	}
}

func mutateOp1(l s.List, where int, m Mutator, vars []string) (s.Sexp, int) {
	head := string(l[0].(s.Atom).Value)
	ml1, r := mutateAt(l[1], where, m, vars)
	return s.List{MkAtom(head), ml1}, r
}

func mutateOp2(l s.List, where int, m Mutator, vars []string) (s.Sexp, int) {
	head := string(l[0].(s.Atom).Value)
	ml1, r1 := mutateAt(l[1], where, m, vars)
	ml2, r2 := mutateAt(l[2], r1, m, vars)
	//fmt.Println("ml1:", ml1, "ml2:", ml2, "r2:", r2)
	return s.List{MkAtom(head), ml1, ml2}, r2
}

func mutateIf0(l s.List, where int, m Mutator, vars []string) (s.Sexp, int) {
	p := l[1]
	zero := l[2]
	nonZero := l[3]
	mp, r1 := mutateAt(p, where, m, vars)
	mzero, r2 := mutateAt(zero, r1, m, vars)
	mnonZero, r3 := mutateAt(nonZero, r2, m, vars)
	return s.List{MkAtom("if0"), mp, mzero, mnonZero}, r3
}

func mutateLambda(l s.List, where int, m Mutator, vars []string) (s.Sexp, int) {
	largs := l[1].(s.List)
	arg1 := string(largs[0].(s.Atom).Value)
	extendedVars := append(vars, arg1)
	if len(largs) == 2 {
		arg2 := string(largs[1].(s.Atom).Value)
		extendedVars = append(extendedVars, arg2)
	}
	body := l[2]
	mbody, r := mutateAt(body, where, m, extendedVars)
	return s.List{MkAtom("lambda"), largs, mbody}, r
}

func mutateFold(l s.List, where int, m Mutator, vars []string) (s.Sexp, int) {
	vec := l[1]
	start := l[2]
	lambda := l[3].(s.List)
	mvec, r1 := mutateAt(vec, where, m, vars)
	mstart, r2 := mutateAt(start, r1, m, vars)
	mlambda, r3 := mutateAt(lambda, r2, m, vars)
	return s.List{MkAtom("fold"), mvec, mstart, mlambda}, r3
}
