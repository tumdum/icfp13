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

func TestEvalProgram(t *testing.T) {
	data := []struct {
		in  string
		env Env
		out uint64
	}{
		{"(lambda (x) (fold x 0 (lambda (y z) (or y z))))", Env{"x": 0x1122334455667788}, 0xff},
		{"(lambda (x) (plus 1 1))", Env{"x": 100}, 2},
	}

	for _, d := range data {
		s := Parse([]byte(d.in))
		if r := EvalProgram(s, d.env["x"]); r != d.out {
			t.Errorf("expected %v, got '%v'", d.out, r)
		}
	}
}

func TestEvalShl1(t *testing.T) {
	data := []struct {
		in  uint64
		out uint64
	}{
		{1, 2},
		{2, 4},
		{3, 6},
		{4, 8},
		{5, 10},
		{6, 12},
		{7, 14},
		{8, 16},
		{9, 18},
		{10, 20},
		{0xffffffff, 0x00000001FFFFFFFE},
	}

	e := "(lambda (x) (shl1 x))"
	for _, d := range data {
		s := Parse([]byte(e))
		if r := Eval(s, Env{"x": d.in}); r != d.out {
			t.Errorf("expected %v, got '%v'", d.out, r)
		}
	}
}

func TestEvalShr1(t *testing.T) {
	data := []struct {
		in  uint64
		out uint64
	}{
		{0x1, 0},
		{0x2, 1},
		{0x3, 1},
		{0x4, 2},
		{0x5, 2},
		{0x6, 3},
		{0x7, 3},
		{0x8, 4},
		{0x9, 4},
		{0x10, 8},
		{0x11, 8},
		{0x12, 9},
		{0xffffffff, 0x000000007FFFFFFF},
	}

	e := "(lambda (x) (shr1 x))"
	for _, d := range data {
		s := Parse([]byte(e))
		if r := Eval(s, Env{"x": d.in}); r != d.out {
			t.Errorf("expected %v, got '%v'", d.out, r)
		}
	}
}

func TestEvalShr4(t *testing.T) {
	data := []struct {
		in  uint64
		out uint64
	}{
		{0x1, 0},
		{0x2, 0},
		{0x3, 0},
		{0x4, 0},
		{0x5, 0},
		{0x6, 0},
		{0x7, 0},
		{0x8, 0},
		{0x9, 0},
		{0x10, 1},
		{0x11, 1},
		{0x12, 1},
		{0xffffffff, 0x000000000FFFFFFF},
	}

	e := "(lambda (x) (shr4 x))"
	for _, d := range data {
		s := Parse([]byte(e))
		if r := Eval(s, Env{"x": d.in}); r != d.out {
			t.Errorf("expected %v, got '%v'", d.out, r)
		}
	}
}

func TestEvalShr16(t *testing.T) {
	data := []struct {
		in  uint64
		out uint64
	}{
		{0x12100000, 0x0000000000001210},
		{0x2, 0},
		{0xffffffff, 0x000000000000FFFF},
	}
	e := "(lambda (x) (shr16 x))"
	for _, d := range data {
		s := Parse([]byte(e))
		if r := Eval(s, Env{"x": d.in}); r != d.out {
			t.Errorf("expected %v, got '%v'", d.out, r)
		}
	}
}

func TestFoldFromTraining(t *testing.T) {
	data := []struct {
		p   string
		in  uint64
		out uint64
	}{
		{"(lambda (x_6545) (fold (shr4 (shr16 x_6545)) x_6545 (lambda (x_6546 x_6547) (or (not x_6547) x_6546))))", 5865335902571436511, 5865335902571436365},
		{"(lambda (x) (fold (shr4 (shr16 x)) x (lambda (a b) (or (not b) a))))", 0x100000000, 0x0000000100000010},
		{"(lambda (x) (fold (shr4 (shr16 x)) x (lambda (a b) (or b a))))", 0x100000000, 0x0000000100000010},
		{"(lambda (x) (not x))", 0x1122334455667788, 0xEEDDCCBBAA998877},
		{"(lambda (x) (fold x x (lambda (a b) (or (not b) a))))", 0x1122334455667788, 0x1122334455667711},
		{"(lambda (x) (fold x 0 (lambda (y z) (or y z))))", 0x1122334455667788, 0x00000000000000FF},
	}

	for _, d := range data {
		s := Parse([]byte(d.p))
		if r := EvalProgram(s, d.in); r != d.out {
			t.Errorf("expected %X, got '%X'", d.out, r)
		}
	}
}

func TestTFoldTraining(t *testing.T) {
	data := []struct {
		p    string
		ins  []uint64
		outs []uint64
	}{
		{"(lambda (x) (fold x 0 (lambda (x a) (or (plus (shl1 a) 1) x))))", []uint64{0, 1, 0x100, 0x1000, 0x10000, 0x100000, 0x1000000, 0x100000000, 0x1122334455667788}, []uint64{0x00000000000000FF, 0x00000000000000FF, 0x00000000000000FF, 0x00000000000004FF, 0x00000000000000FF, 0x00000000000002FF, 0x00000000000000FF, 0x00000000000000FF, 0x0000000000005FFF}},
		{"(lambda (x) (fold x 0 (lambda (x a) (plus x a))))", []uint64{0, 1, 0x100, 0x1000, 0x10000, 0x100000, 0x1000000, 0x100000000, 0x1122334455667788}, []uint64{0x0000000000000000, 0x0000000000000001, 0x0000000000000001, 0x0000000000000010, 0x0000000000000001, 0x0000000000000010, 0x0000000000000001, 0x0000000000000001, 0x0000000000000264}},
	}

	for _, d := range data {
		e := Parse([]byte(d.p))
		for idx, in := range d.ins {
			if r := EvalProgram(e, in); r != d.outs[idx] {
				t.Errorf("%d: form %X(%d) expected '%X(%d)', got '%X(%d)'", idx, in, in, d.outs[idx], d.outs[idx], r, r)
			}
		}
	}
}
