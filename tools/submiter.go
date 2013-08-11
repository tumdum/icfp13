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
  "runtime"
)

const TestDataSize = 70

func RandomInput(size int) []string {
	ret := make([]string, size)
	for i := 0; i < size; i++ {
		v1 := uint64(rand.Int63())
		v2 := uint64(rand.Int63())
		v := v1 + v2
		ret[i] = "0x" + strconv.FormatUint(v, 16)
	}
ret = append(ret, "0xFFFFFFFFF0000001")
ret = append(ret, "0x0000000000000040")
ret = append(ret, "0xE108108020400180")
ret = append(ret, "0x4001802001020200")
ret = append(ret, "0xFFFF0000FFFF0001")
ret = append(ret, "0x8707040005020503")
ret = append(ret, "0xFFFFFFFFFFFE0000")
ret = append(ret, "0x0000000000000030")
ret = append(ret, "0xF0F0F0F0F0F0F0F1")
ret = append(ret, "0x00000FFFFFFFFFFE")
ret = append(ret, "0x3333333333333333")
ret = append(ret, "0x8000000000000004")
ret = append(ret, "0x0000000000000004")
ret = append(ret, "0x0200FE062020201E")
ret = append(ret, "0xAAAAAAAAAAAAAAAB")
ret = append(ret, "0x80000000D652A6AA")
ret = append(ret, "0x8000000000000000")
ret = append(ret, "0x0161C08342C010DF")
	ret = append(ret, "0x0")
	ret = append(ret, "0xffffffffffffffff")
ret = append(ret, "0x0000000900000000")
  ret = append(ret, "0x0000000000000020")
  ret = append(ret, "0x0000000000000001")
  ret = append(ret, "0x0000000000000010")
  ret = append(ret, "0x000000000010FFFF")
  ret = append(ret, "0x00000000FFE00000")
  ret = append(ret, "0x7FFFFFFFFFFFFFFF")
  ret = append(ret, "0xFFFFF00000FFFFF0")
  ret = append(ret, "0x0000000000010008")
	ret = append(ret, "0x8000000000008000")
	return ret
}

func GetConstraints(id string) []icfp13.Constraint {
	ri := RandomInput(TestDataSize)
	ereq := service.EvalRequest{id, "", ri}
  var eresp *service.EvalResponse
  var e error
  for {
	eresp, e = service.Eval(ereq)
  if e != nil && strings.HasPrefix(e.Error(), "429") {
			time.Sleep(5 * time.Second)
			continue
		}

	if e != nil {
		panic(e)
	}
  break
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
  ret := make(chan icfp13.Solution)
	go icfp13.FindProgramPar(constraints, ops, size, ret)
  solution := <-ret
	fmt.Println(solution)
	solstr := solution.Prog.String()
	solstr = strings.Replace(solstr, " __0", " 0", -1)
	solstr = strings.Replace(solstr, " __1", " 1", -1)
  solstr = strings.Replace(solstr, "const_","", -1)

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
    if gs.Status == "mismatch"  {
      os.Exit(1)
    }
    if gs.Status == "error" {
      os.Exit(1)
    }

    break
	}
}

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())
	id := os.Args[1]
	size, _ := strconv.Atoi(os.Args[2])
	ops := make([]string, 0)
	for i := 3; i < len(os.Args); i++ {
		ops = append(ops, os.Args[i])
	}

	fmt.Printf("Will try to solve program with id '%v', size '%v' and ops '%v'.\n", id, size, ops)
	Solve(id, size, ops)
}
