package main

import (
	"bitbucket.org/tumdum/icfp13/service"
	"fmt"
)

func main() {
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
}
