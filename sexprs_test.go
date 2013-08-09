package icfp13

import "testing"

func TestEvalOp2(t *testing.T) {
	data := []struct {
		input  string
		output uint64
	}{
		{"(or 0 1)", 1},
		{"(or 1 0)", 1},
		{"(or 0 0)", 0},
		{"(or 1 1)", 1},
		{"(or (or (or 1 1) 1) 1)", 1},
		{"(and 0 0)", 0},
		{"(and 0 1)", 0},
		{"(and 1 0)", 0},
		{"(and 1 1)", 1},
		{"(xor 0 0)", 0},
		{"(xor 0 1)", 1},
		{"(xor 1 0)", 1},
		{"(xor 1 1)", 0},
		{"(plus 1 (plus 1 1))", 3},
	}
	for _, d := range data {
		s := Parse([]byte(d.input))
		if r := Eval(s, make(Env)); r != d.output {
			t.Errorf("expected %v, got '%v'", d.output, r)
		}
	}
}

func TestEvalId(t *testing.T) {
	data := []struct {
		expr string
		env  Env
		out  uint64
	}{
		{"(plus x y)", Env{"x": 55, "y": 45}, 100},
	}

	for _, d := range data {
		s := Parse([]byte(d.expr))
		if r := Eval(s, d.env); r != d.out {
			t.Errorf("expected %v, got '%v'", d.out, r)
		}
	}
}

func TestEvalOp1(t *testing.T) {
	data := []struct {
		in  string
		out uint64
	}{
		{"(not 0)", 0xffffffffffffffff},
		{"(not (not 0))", 0},
		{"(not (not 1))", 1},
	}
	for _, d := range data {
		s := Parse([]byte(d.in))
		if r := Eval(s, make(Env)); r != d.out {
			t.Errorf("expected %v, got '%v'", d.out, r)
		}
	}
}

func TestEvalIf0(t *testing.T) {
	data := []struct {
		in  string
		env Env
		out uint64
	}{
		{"(if0 0 0 1)", make(Env), 0},
		{"(if0 1 0 1)", make(Env), 1},
		{"(if0 x y z)", Env{"x": 0, "y": 100, "z": 200}, 100},
	}

	for _, d := range data {
		s := Parse([]byte(d.in))
		if r := Eval(s, d.env); r != d.out {
			t.Errorf("expected %v, got '%v'", d.out, r)
		}
	}
}

func TestEvalLambda(t *testing.T) {
	data := []struct {
		in  string
		env Env
		out uint64
	}{
		{"(lambda (x y) (plus x y))", Env{"x": 100, "y": 50}, 150},
		{"(lambda (x y) (plus (plus x x) y))", Env{"x": 100, "y": 3}, 203},
	}

	for _, d := range data {
		s := Parse([]byte(d.in))
		if r := Eval(s, d.env); r != d.out {
			t.Errorf("expected %v, got '%v'", d.out, r)
		}
	}
}

func TestEvalFold(t *testing.T) {
	data := []struct {
		in  string
		env Env
		out uint64
	}{
		{"(fold x 0 (lambda (y z) (or y z)))", Env{"x": 0x1122334455667788}, 0xff},
	}

	for _, d := range data {
		s := Parse([]byte(d.in))
		if r := Eval(s, d.env); r != d.out {
			t.Errorf("expected %v, got '%v'", d.out, r)
		}
	}
}
