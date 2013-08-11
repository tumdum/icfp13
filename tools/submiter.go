package main

import (
	"bitbucket.org/tumdum/icfp13"
	"bitbucket.org/tumdum/icfp13/service"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const TestDataSize = 50

func RandomInput(size int) []string {
	ret := make([]string, size)
	for i := 0; i < size; i++ {
		v1 := uint64(rand.Int63())
		v2 := uint64(rand.Int63())
		v := v1 + v2
		ret[i] = "0x" + strconv.FormatUint(v, 16)
	}
	ret = append(ret, "0x0")
	ret = append(ret, "0xffffffffffffffff")
	return ret
}

func GetConstraints(id string) []icfp13.Constraint {
	ri := RandomInput(TestDataSize)
	ereq := service.EvalRequest{id, "", ri}
	eresp, e := service.Eval(ereq)
	if e != nil {
		panic(e)
	}

	cons := make([]icfp13.Constraint, len(ri))
	for i := 0; i < len(cons); i++ {
		in, _ := strconv.ParseUint(ri[i], 0, 64)
		out, _ := strconv.ParseUint(eresp.Outputs[i], 0, 64)
		cons[i] = icfp13.Constraint{in, out}
	}
	return cons
}

func Solve(id string, size int, ops []string) {
	constraints := GetConstraints(id)
	solution := icfp13.FindProgramPar(constraints, ops, size)
	fmt.Println(solution)
	solstr := solution.String()
	solstr = strings.Replace(solstr, " __0", " 0", -1)
	solstr = strings.Replace(solstr, " __1", " 1", -1)

	for {

		gr := service.GuessRequest{id, solstr}
		fmt.Println(gr)
		gs, err := service.Guess(gr)
		if err != nil && strings.HasPrefix(err.Error(), "429") {
			time.Sleep(5 * time.Second)
			continue
		}
		if err != nil {
			fmt.Println("!!! ERROR: ", err)
			break
		}
		fmt.Println(gs)
    break
	}
}

func main() {
	id := os.Args[1]
	size, _ := strconv.Atoi(os.Args[2])
	ops := make([]string, 0)
	for i := 3; i < len(os.Args); i++ {
		ops = append(ops, os.Args[i])
	}

	fmt.Printf("Will try to solve program with id '%v', size '%v' and ops '%v'.\n", id, size, ops)
	Solve(id, size, ops)
}
