package icfp13

import (
	"fmt"
	s "github.com/eadmund/sexprs"
	"math/rand"
	"runtime"
	"sort"
)

const StartSexp = "(lambda (x) x)"
const Percent = 0.6
const MaxGenerationSize = 10000

type Constraint struct {
	In, Out uint64
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
	/*if prefixSize >= 2*MaxGenerationSize {
	  fmt.Println("hard trim!")
	  prefixSize = 2*MaxGenerationSize
	}*/
	return sols[:prefixSize]
}

type NextGenReq struct {
	start       s.Sexp
	constraints []Constraint
	ops         []string
}

func Generator(reqs chan NextGenReq, out chan Solutions, stop chan bool, size int) {
	for {
		select {
		case req := <-reqs:
			s := NextGeneration(req.start, req.constraints, req.ops, size)
			out <- s
		case <-stop:
			break
		}
	}
}

func CollectResults(out chan Solutions, count int, ret chan Solutions, solret chan Solution) {
	newsols := make(Solutions, 0)
	for i := 0; i < count; i++ {
		ss := <-out
		newsols = append(newsols, ss...)
		for _, sol := range ss {
			if sol.score == 1.0 {
				// fmt.Println("found solution",sol.prog)
				solret <- sol
			}
		}
	}
	ret <- newsols
}

func RemoveTooBig(sols Solutions, targetSize int) Solutions {
	r := make(Solutions, 0)
	dropped := 0
	for _, sol := range sols {
		if Size(sol.Prog) < targetSize /*|| (Size(sol.Prog) == targetSize && sol.score == 1.0)*/ {
			r = append(r, sol)
		} else {
			dropped++
		}
	}
	fmt.Println("Dropped", dropped, "too big solutions")
	return r
}

func GenSize(iterNumber int) int {
	if iterNumber == 0 || iterNumber == 1 {
		return 4 * NewGenerationSize
	}
	return NewGenerationSize
}

func HasBonus(ops []string) bool {
	for _, op := range ops {
		if op == "bonus" {
			return true
		}
	}
	return false
}

func HasTfold(ops []string) bool {
	for _, op := range ops {
		if op == "tfold" {
			return true
		}
	}
	return false
}

func FindProgramPar(constraints []Constraint, ops []string, size int, solret chan Solution) s.Sexp {
	req := make(chan NextGenReq, 10)
	out := make(chan Solutions)
	merged := make(chan Solutions)
	stop := make(chan bool)
	for i := 0; i < runtime.NumCPU(); i++ {
		go Generator(req, out, stop, NewGenerationSize)
	}
	/* go Generator(req, out, stop, NewGenerationSize)
	go Generator(req, out, stop, NewGenerationSize)
	go Generator(req, out, stop, NewGenerationSize)
	*/

	var start s.Sexp
	var solution s.Sexp
	if HasTfold(ops) {
		start = Parse([]byte("(lambda (const_x) (fold const_x const_0 (lambda (const_x y) e)))"))
	} else if HasBonus(ops) {
		start = Parse([]byte("(lambda (x) (if0 x x x))"))
		newOps := make([]string, 0)
		for _, op := range ops {
			if op != "bonus" && op != "if0" {
				newOps = append(newOps, op)
			}
		}
		ops = newOps
	} else {
		start = Parse([]byte(StartSexp))
	}
	sols := NextGeneration(start, constraints, ops, 10*NewGenerationSize)
	i := 0
	lastBestScore := 0.0
	for {
		fmt.Println("iter:", i, len(sols), sols[0].score)
		newsols := make(Solutions, 0)
		go CollectResults(out, len(sols), merged, solret)
		for _, sol := range sols {
			req <- NextGenReq{sol.Prog, constraints, ops}
		}

		newsols = <-merged

		sort.Sort(newsols)
		sols = TakeBestPercent(Percent, newsols)
		if sols[0].score > 0.0 {
			sols = RemoveTooBig(sols, size)
			if len(sols) == 0 || lastBestScore >= sols[0].score {
				sols = TakeBestPercent(Percent, NextGeneration(start, constraints, ops, 10*NewGenerationSize))
				lastBestScore = 0.0
				i = 0
			}
			if sols[0].score == 0.0 && len(sols) >= MaxGenerationSize*2 {
				i = 0
				lastBestScore = 0
				sols = NextGeneration(start, constraints, ops, 10*NewGenerationSize)
			}
		}
		i++
	}
	return solution
}

type Solution struct {
	Prog  s.Sexp
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

func NextGeneration(e s.Sexp, constraints []Constraint, ops []string, size int) Solutions {
	nextGeneration := GenNewRandomUniqGenerationUsing(e, ops, size)
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
		r := EvalProgram(e, con.In)
		if r == con.Out {
			ok = ok + 1.0
		}
	}
	return ok / total
}
