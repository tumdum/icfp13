package icfp13

import (
	"bytes"
	"github.com/eadmund/sexprs"
)

type Env map[string]uint64

func Parse(input []byte) sexprs.Sexp {
	input = bytes.Replace(input, []byte(" 0"), []byte(" __0"), -1)
	input = bytes.Replace(input, []byte(" 1"), []byte(" __1"), -1)
	s, r, e := sexprs.Parse(input)
	if len(r) != 0 {
		panic("rest not empty: " + string(r))
	}
	if e != nil {
		panic("failed to parse")
	}
	return s
}

func Eval(e sexprs.Sexp, input Env) uint64 {
	switch e := e.(type) {
	case sexprs.List:
		return evalList(e, input)
	case sexprs.Atom:
		return evalAtom(e, input)
	}
	return 4
}

func evalList(l sexprs.List, input Env) uint64 {
	head := string(l[0].(sexprs.Atom).Value)
	switch head {
	case "or":
		return evalOr(l[1], l[2], input)
	case "and":
		return evalAnd(l[1], l[2], input)
	case "xor":
		return evalXor(l[1], l[2], input)
	case "plus":
		return evalPlus(l[1], l[2], input)
	case "not":
		return evalNot(l[1], input)
	case "if0":
		return evalIf0(l[1], l[2], l[3], input)
	case "lambda":
		return evalLambda(l[1], l[2], input)
	default:
		panic("unknown list head: " + head)
	}
}

func evalOr(e1, e2 sexprs.Sexp, input Env) uint64 {
	e1v := Eval(e1, input)
	e2v := Eval(e2, input)
	return e1v | e2v
}

func evalAnd(e1, e2 sexprs.Sexp, input Env) uint64 {
	e1v := Eval(e1, input)
	e2v := Eval(e2, input)
	return e1v & e2v
}

func evalXor(e1, e2 sexprs.Sexp, input Env) uint64 {
	e1v := Eval(e1, input)
	e2v := Eval(e2, input)
	return e1v ^ e2v
}

func evalPlus(e1, e2 sexprs.Sexp, input Env) uint64 {
	e1v := Eval(e1, input)
	e2v := Eval(e2, input)
	return e1v + e2v
}

func evalNot(e sexprs.Sexp, input Env) uint64 {
	return ^Eval(e, input)
}

func evalIf0(p, zero, nonZero sexprs.Sexp, input Env) uint64 {
	pv := Eval(p, input)
	if pv == 0 {
		return Eval(zero, input)
	} else {
		return Eval(nonZero, input)
	}
}

func evalLambda(params, body sexprs.Sexp, input Env) uint64 {
	/* p1 and p2 should already be in input!
	pList := params.(sexprs.List)
	p1 := string(pList[0].(sexprs.Atom).Value)
	p2 := string(pList[1].(sexprs.Atom).Value)
	p1v := input[p1]
	p2v := input[p2]
	*/
	return Eval(body, input)
}

func evalAtom(e sexprs.Atom, input Env) uint64 {
	es := string(e.Value)
	switch es {
	case "__0":
		return 0
	case "__1":
		return 1
	default:
		return input[es]
	}
}
