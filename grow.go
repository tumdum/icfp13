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

type Mutator func(s.Sexp) s.Sexp

func MutateAt(e s.Sexp, where int, m Mutator) (s.Sexp, error) {
	// fmt.Println("---------")
	me, n := mutateAt(e, where, m)
	if n != 0 {
		return nil, errors.New("failed to mutate")
	}
	return me, nil
}

func mutateAt(e s.Sexp, where int, m Mutator) (s.Sexp, int) {
	switch e := e.(type) {
	case s.List:
		return mutateList(e, where, m)
	case s.Atom:
		return mutateAtom(e, where, m)
	default:
		panic("mutate unknown type")
	}
}

func mutateAtom(a s.Atom, where int, m Mutator) (s.Sexp, int) {
	// fmt.Println("mutateAtom", a, where)
	if where == 1 {
		return m(a), 0
	}
	if where > 0 {
		where--
	}
	return MkAtom(string(a.Value)), where
}

func mutateList(l s.List, where int, m Mutator) (s.Sexp, int) {
	head := string(l[0].(s.Atom).Value)
	switch head {
	case "not", "shl1", "shr1", "shr4", "shr16":
		return mutateOp1(l, where, m)
	case "or", "and", "xor", "plus":
		return mutateOp2(l, where, m)
	default:
		panic("mutate list with unknown head: " + head)
	}
}

func mutateOp1(l s.List, where int, m Mutator) (s.Sexp, int) {
	head := string(l[0].(s.Atom).Value)
	ml1, r := mutateAt(l[1], where, m)
	return s.List{MkAtom(head), ml1}, r
}

func mutateOp2(l s.List, where int, m Mutator) (s.Sexp, int) {
	head := string(l[0].(s.Atom).Value)
	ml1, r1 := mutateAt(l[1], where, m)
	ml2, r2 := mutateAt(l[2], r1, m)
	//fmt.Println("ml1:", ml1, "ml2:", ml2, "r2:", r2)
	return s.List{MkAtom(head), ml1, ml2}, r2
	/*if where == 1 {
	    return s.List{ MkAtom(head), m(l[1]), l[2] }, 0
	  } else if where == 2 {
	    return s.List{ MkAtom(head), l[1], m(l[2]) }, 0
	  }
	  return s.List{ MkAtom(head), l[1], l[2] }, where*/
}
