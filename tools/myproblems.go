package main

import (
  "bitbucket.org/tumdum/icfp13/service"
  "sort"
  "fmt"
)

type Problems []service.Problem

func (p Problems) Len() int {
  return len(p)
}

func (p Problems) Less(i, j int) bool {
  return p[i].Size > p[j].Size
}

func (p Problems) Swap(i, j int) {
  p[i], p[j] = p[j], p[i]
}

func main() {
  problems, err := service.MyProblems()
  if err != nil {
    panic(err)
  }
  problemsToSort := Problems(problems)
  sort.Sort(problemsToSort)
  for _, problem := range problemsToSort {
    if !problem.Solved && problem.TimeLeft == nil {
      fmt.Printf("timelimit -t 301 ./submiter %v %v ", problem.Id,problem.Size)
      for _, op := range problem.Operators {
        fmt.Printf("%v ", op)
      }
      fmt.Println()
    }
  }
}
