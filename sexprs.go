package icfp13

import (
	"bytes"
	"github.com/eadmund/sexprs"
)

type Env map[string]uint64

var listFuns map[string]func(sexprs.List, Env) uint64

func init() {
	listFuns = make(map[string]func(sexprs.List, Env) uint64)
	listFuns["or"] = evalOr
	listFuns["and"] = evalAnd
	listFuns["xor"] = evalXor
	listFuns["plus"] = evalPlus
	listFuns["not"] = evalNot
	listFuns["shl1"] = evalShl1
	listFuns["shr1"] = evalShr1
	listFuns["shr4"] = evalShr4
	listFuns["shr16"] = evalShr16
	listFuns["if0"] = evalIf0
	listFuns["lambda"] = evalLambda
	listFuns["fold"] = evalFold
}

func Parse(input []byte) sexprs.Sexp {
	input = bytes.Replace(input, []byte(" 0"), []byte(" __0"), -1)
	input = bytes.Replace(input, []byte(" 1"), []byte(" __1"), -1)
	s, r, e := sexprs.Parse(input)
	if len(r) != 0 {
		panic("rest not empty: " + string(r))
	}
	if e != nil {
		panic("failed to parse: " + string(input))
	}
	return s
}

func EvalProgram(p sexprs.Sexp, input uint64) uint64 {
	arg := string(p.(sexprs.List)[1].(sexprs.List)[0].(sexprs.Atom).Value)
	r := Eval(p, Env{arg: input})
	return r
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
	return listFuns[head](l, input)
	/*switch head {
		case "or":
			// return evalOr(l[1], l[2], input)
	    return evalOr(l, input)
		case "and":
			// return evalAnd(l[1], l[2], input)
	    return evalAnd(l, input)
		case "xor":
			return evalXor(l, input)
		case "plus":
			return evalPlus(l, input)
		case "not":
			return evalNot(l, input)
		case "shl1":
			return evalShl1(l, input)
		case "shr1":
			return evalShr1(l, input)
		case "shr4":
			return evalShr4(l, input)
		case "shr16":
			return evalShr16(l, input)
		case "if0":
			return evalIf0(l, input)
		case "lambda":
			return evalLambda(l, input)
		case "fold":
			return evalFold(l, input)
		default:
			panic("unknown list head: " + head)
		}*/
}

func evalOr(l sexprs.List, input Env) uint64 {
	e1v := Eval(l[1], input)
	e2v := Eval(l[2], input)
	return e1v | e2v
}

func evalAnd(l sexprs.List, input Env) uint64 {
	e1v := Eval(l[1], input)
	e2v := Eval(l[2], input)
	return e1v & e2v
}

func evalXor(l sexprs.List, input Env) uint64 {
	e1v := Eval(l[1], input)
	e2v := Eval(l[2], input)
	return e1v ^ e2v
}

func evalPlus(l sexprs.List, input Env) uint64 {
	e1v := Eval(l[1], input)
	e2v := Eval(l[2], input)
	return e1v + e2v
}

func evalNot(l sexprs.List, input Env) uint64 {
	return ^Eval(l[1], input)
}

func evalShl1(l sexprs.List, input Env) uint64 {
	return Eval(l[1], input) << 1
}

func evalShr1(l sexprs.List, input Env) uint64 {
	return Eval(l[1], input) >> 1
}

func evalShr4(l sexprs.List, input Env) uint64 {
	return Eval(l[1], input) >> 4
}

func evalShr16(l sexprs.List, input Env) uint64 {
	return Eval(l[1], input) >> 16
}

// func evalIf0(p, zero, nonZero sexprs.Sexp, input Env) uint64 {
func evalIf0(l sexprs.List, input Env) uint64 {
	pv := Eval(l[1], input)
	if pv == 0 {
		return Eval(l[2], input)
	} else {
		return Eval(l[3], input)
	}
}

func evalLambda(l sexprs.List, input Env) uint64 {
	/* p1 and p2 should already be in input!
	pList := params.(sexprs.List)
	p1 := string(pList[0].(sexprs.Atom).Value)
	p2 := string(pList[1].(sexprs.Atom).Value)
	p1v := input[p1]
	p2v := input[p2]
	*/
	return Eval(l[2], input)
}

func evalFold(l sexprs.List, input Env) uint64 {
	vecv := Eval(l[1], input)
	acc := Eval(l[2], input)
	for i := 0; i <= 7; i++ {
		left := (vecv << uint((7-i)*8)) >> uint(7*8)
		localEnv := make(Env)
		for k, v := range input {
			localEnv[k] = v
		}
		lparams := l[3].(sexprs.List)[1].(sexprs.List)
		x := string(lparams[0].(sexprs.Atom).Value)
		y := string(lparams[1].(sexprs.Atom).Value)
		localEnv[x] = left
		localEnv[y] = acc
		acc = Eval(l[3], localEnv)
	}
	return acc
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
