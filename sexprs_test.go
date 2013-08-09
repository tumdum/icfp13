package icfp13

import "testing"

func TestEvalOp2(t *testing.T) {
  data := []struct {
    input string
    output uint64
  } {
    { "(or 0 1)", 1 },
    { "(or 1 0)", 1 },
    { "(or 0 0)", 0 },
    { "(or 1 1)", 1 },
    { "(or (or (or 1 1) 1) 1)", 1 },
    { "(and 0 0)", 0},
    { "(and 0 1)", 0},
    { "(and 1 0)", 0},
    { "(and 1 1)", 1},
    { "(xor 0 0)", 0},
    { "(xor 0 1)", 1},
    { "(xor 1 0)", 1},
    { "(xor 1 1)", 0},
    { "(plus 1 (plus 1 1))", 3},
  }
  for _,d := range data {
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
  } {
    { "(plus x y)", Env{"x":55,"y":45}, 100},
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
    in string
    out uint64
  } {
    { "(not 0)", 0xffffffffffffffff },
    { "(not (not 0))", 0 },
    { "(not (not 1))", 1 },
  }
  for _, d := range data {
    s := Parse([]byte(d.in))
    if r := Eval(s, make(Env)); r != d.out {
      t.Errorf("expected %v, got '%v'", d.out, r)
    }
  }
}
