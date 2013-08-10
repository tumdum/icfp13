package main

import (
  "bitbucket.org/tumdum/icfp13"
  "runtime"
  "fmt"
)

func run(p string, ops []string, size int) {
  e := icfp13.Parse([]byte(p))
  fmt.Println("Starting work on solution for",p,"(size ",icfp13.Size(e),")")
  cons := icfp13.GenConstrains(e, 100)
  icfp13.FindProgramPar(cons, ops, size)
}

func main() {
  runtime.GOMAXPROCS(4)
  p := "(lambda (x) (plus 1 1))"
  run(p,[]string{"plus"},5)
  run("(lambda (x) (or (plus x 1) (shr4 x)))",[]string{"shr4","plus","or"},10)
  run("(lambda (x_5856) (and (shr4 x_5856) (shr1 (shr1 x_5856))))", []string{ "and", "shr1", "shr4" },7) // size 7
  run("(lambda (x_7325) (and (not (plus 1 (shr16 x_7325))) x_7325))", []string{ "and", "not", "plus", "shr16" },8) // size 8
  run("(lambda (x_7942) (if0 (and (and x_7942 1) 1) 0 x_7942))", []string{"and","if0"},9) // size 9
  run("(lambda (x_11238) (shr1 (if0 (or (shr1 x_11238) 1) (shr1 (not x_11238)) x_11238)))", []string{"if0", "not", "or", "shr1"},11) // size 11
  run("(lambda (x_17733) (xor (not (shr1 (plus (if0 (shr1 (xor (not x_17733) 0)) x_17733 1) x_17733))) x_17733))", []string{ "if0", "not", "plus", "shr1", "xor" },15) // size 15
  // long:
  run("(lambda (x_9921) (fold x_9921 0 (lambda (x_9921 x_9922) (or (xor (shr1 x_9921) 0) 1))))", []string{"or", "shr1", "tfold", "xor" }, 11)
}
