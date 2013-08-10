package main

import "bitbucket.org/tumdum/icfp13"

func run(p string) {
  e := icfp13.Parse([]byte(p))
  cons := icfp13.GenConstrains(e, 1000)
  icfp13.FindProgram(cons)

  // icfp13.MutateAt(e, 1, icfp13.GenSexp)
}

func main() {
  p := "(lambda (x) (not x))"
  run(p)
}
