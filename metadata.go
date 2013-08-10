package icfp13

import (
	s "github.com/eadmund/sexprs"
)

func Size(e s.Sexp) int {
	switch e := e.(type) {
	case s.Atom:
		return atomSize(e)
	case s.List:
		return listSize(e)
	default:
		panic("here be dragons")
	}
}

func atomSize(a s.Atom) int {
	return 1
}

func listSize(l s.List) int {
	head := string(l[0].(s.Atom).Value)
	switch head {
	case "if0":
		return 1 + Size(l[1]) + Size(l[2]) + Size(l[3])
	case "fold":
		return foldSize(l[1], l[2], l[3])
	case "lambda":
		return 1 + Size(l[2])
	case "or", "and", "xor", "plus":
		return 1 + Size(l[1]) + Size(l[2])
	case "not", "shl1", "shr1", "shr4", "shr16":
		return 1 + Size(l[1])
	default:
		panic("?!")
	}
}

func foldSize(vec, start, lambda s.Sexp) int {
	lbody := lambda.(s.List)[2]
	return 2 + Size(vec) + Size(start) + Size(lbody)
}
