package icfp13

import (
	"fmt"
	s "github.com/eadmund/sexprs"
	"math/rand"
	"sort"
)

const StartSexp = "(lambda (x) x)"
const Percent = 0.7
const MaxGenerationSize = 1000

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

func TakeBestPercent(percent float32, sols Solutions) Solutions {
	l := sols.Len()
	prefixSize := int(percent*float32(l)) + 1
	if prefixSize > MaxGenerationSize && sols[0].score > 0 {
		fmt.Println("trim!")
		prefixSize = MaxGenerationSize
	}
	return sols[:prefixSize]
}

func FindProgram(constraints []Constraint, ops []string) {
	start := Parse([]byte(StartSexp))
	sols := TakeBestPercent(Percent, NextGeneration(start, constraints, ops))
	i := 0
	for {
		fmt.Println("iter:", i, len(sols))
		newsols := make(Solutions, 0)
		for _, sol := range sols {
			newsols = append(newsols, NextGeneration(sol.prog, constraints, ops)...)
		}
		sols = TakeBestPercent(Percent, append(sols, newsols...))
		sort.Sort(sols)
		if sols[0].score == 1.0 {
			fmt.Println("found solution:", sols[0].prog)
			Score(sols[0].prog, constraints)
			break
		} else {
			fmt.Println("best score:", sols[0].score)
		}
		i++
	}
}

type Solution struct {
	prog  s.Sexp
	score float64
}

type Solutions []Solution

func (s Solutions) Len() int {
	return len(s)
}

func (s Solutions) Less(i, j int) bool {
	return s[i].score > s[j].score
}

func (s Solutions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func NextGeneration(e s.Sexp, constraints []Constraint, ops []string) Solutions {
	nextGeneration := GenNewRandomUniqGenerationUsing(e, ops)
	sols := make([]Solution, 0)
	for _, next := range nextGeneration {
		score := Score(next, constraints)
		sols = append(sols, Solution{next, score})
	}
	return sols
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
