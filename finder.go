package icfp13

import (
	s "github.com/eadmund/sexprs"
  "fmt"
  "math/rand"
)

const StartSexp = "(lambda (x) x)"

type Constraint struct {
  in, out uint64
}

func Uint64() uint64 {
  v1 := uint64(rand.Int63())
  v2 := uint64(rand.Int63())
  return v1 + v2
}

func GenConstrains(e s.Sexp, s int) []Constraint {
  ret := make([]Constraint, s)
  for i := 0; i < s; i++ {
    input := Uint64()
    value := EvalProgram(e, input)
    ret[i] = Constraint{input, value}
  }
  return ret
}

func FindProgram(constraints []Constraint) {
  start := Parse([]byte(StartSexp))
  nextGeneration := GenNewRandomUniqGeneration(start)
  for _, next := range nextGeneration {
    score := Score(next, constraints)
    fmt.Println(score, next)
  }
}

func Score(e s.Sexp, cons []Constraint) float64 {
  total := float64(len(cons))
  ok := float64(0)
  for _, con := range cons {
    r := EvalProgram(e, con.in)
    if r == con.out {
      ok = ok + 1.0
    }
  }
  return ok / total
}

