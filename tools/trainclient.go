package main

import (
	"bitbucket.org/tumdum/icfp13"
	"bitbucket.org/tumdum/icfp13/service"
	"fmt"
	"math/rand"
	"strconv"
)

func RandomInput(size int) []string {
	ret := make([]string, size)
	for i := 0; i < size; i++ {
		v1 := uint64(rand.Int63())
		v2 := uint64(rand.Int63())
		v := v1 + v2
		ret[i] = "0x" + strconv.FormatUint(v, 16)
	}
	return ret
}

func CheckEval(psize int) {
	for {
		ri := RandomInput(255)
		problem, e := service.Train(service.TrainRequest{psize, []string{"tfold"}})
		if e != nil {
			continue
		}
		fmt.Println("problem:", problem)
		ereq := service.EvalRequest{problem.Id, problem.Challenge, ri}
		eresp, e := service.Eval(ereq)
		if e != nil {
			continue
		}
		// fmt.Println("eresp:",eresp)
		prog := icfp13.Parse([]byte(problem.Challenge))
		for i, ins := range ri {
			in, _ := strconv.ParseUint(ins, 0, 64)
			progRet := icfp13.EvalProgram(prog, in)
			out, _ := strconv.ParseUint(eresp.Outputs[i], 0, 64)
			if progRet != out {
				fmt.Printf("for p = \"%v\", and arg = %d\n%d != %d\n", problem.Challenge, in, progRet, out)
			}
		}
	}
}

func verify(p string) {
  inputs := []string{}
  for i := uint64(0); i < 255; i++ {
    // inputs = append(inputs, "0x" + strconv.FormatUint(i*10, 16))
    inputs = append(inputs, "0x100")
  }
  prog := icfp13.Parse([]byte(p))
  ereq := service.EvalRequest{"", p, inputs}
  eresp, _ := service.Eval(ereq)
  fmt.Println(eresp)
  for i, val := range eresp.Outputs {
    in, _ := strconv.ParseUint(inputs[i], 0, 64)
    progRet := icfp13.EvalProgram(prog, in)
    out, _ := strconv.ParseUint(val, 0, 64)
    fmt.Printf("%d: %s => %X <= %X\n", i, inputs[i], out, progRet)
  }
}

func main() {
	CheckEval(18)
  // verify("(lambda (x) (fold x 0 (lambda (x a) (plus x a))))")
}

func Xmain() {
	problem, e := service.Train(service.TrainRequest{4, []string{}})
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("problem:", problem)
	ereq := service.EvalRequest{problem.Id, problem.Challenge, []string{"0", "1", "2"}}
	eresp, e := service.Eval(ereq)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(eresp)
	fmt.Println(RandomInput(10))
}
