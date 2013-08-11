package main

import (
"bitbucket.org/tumdum/icfp13"
"os"
"runtime/pprof"
)

func main() {
  f, err := os.Create("profiling.prof")
        if err != nil {
          panic(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()

  p := "(lambda (x_80999) (fold x_80999 0 (lambda (x_80999 x_81000) (xor x_80999 (shr1 (not (or (shr1 (or (xor x_80999 x_80999) (and x_80999 (not (and (if0 (plus (shr16 x_80999) (and 0 x_81000)) 1 x_80999) x_81000))))) x_81000)))))))"
  prog := icfp13.Parse([]byte(p))
  for i := 0; i < 100000; i++ {
    icfp13.EvalProgram(prog, 1024)
    icfp13.EvalProgram(prog, 213531479)
  }
}
