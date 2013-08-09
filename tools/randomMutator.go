package main

import (
	"bitbucket.org/tumdum/icfp13"
	"fmt"
	"math/rand"
)

func main() {
	str := "(lambda (x_38078) (fold (if0 (shr1 (shr1 (xor x_38078 x_38078))) x_38078 x_38078) 0 (lambda (x_38079 x_38080) (not (or x_38079 x_38080)))))"
	e := icfp13.Parse([]byte(str))
	max := icfp13.CountMutationPoints(e)
	for i := 0; i < 1000; i++ {
		idx := rand.Int() % max
		me, err := icfp13.MutateAt(e, idx+1, icfp13.GenSexp)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println(me)
		}
	}
}
