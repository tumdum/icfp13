package main

import (
	"bitbucket.org/tumdum/icfp13"
	"fmt"
)

func main() {
	str := "(lambda (x_38078) (fold (if0 (shr1 (shr1 (xor x_38078 x_38078))) x_38078 x_38078) 0 (lambda (x_38079 x_38080) (not (or x_38079 x_38080)))))"
	str = "(lambda (x_29213) (fold (or (shl1 (or x_29213 x_29213)) 0) x_29213 (lambda (x_29214 x_29215) (if0 x_29215 x_29215 x_29214))))"
	str = "(lambda (x_10085) (fold x_10085 0 (lambda (x_10085 x_10086) (and (if0 x_10085 0 1) x_10085))))"
	str = "(lambda (x_36848) (fold (if0 (shr16 (or (not x_36848) x_36848)) x_36848 x_36848) 0 (lambda (x_36849 x_36850) (if0 x_36850 x_36849 x_36849))))"
	e := icfp13.Parse([]byte(str))
	for i := uint64(0); i < 10; i++ {
		v := icfp13.Eval(e, icfp13.Env{"x_36848": i})
		fmt.Printf("%v ==> %X\n", i, v)
	}
	fmt.Printf(" %X\n", icfp13.Eval(e, icfp13.Env{"x_36848": 0xffffffffffffffff}))

	/*max := icfp13.CountMutationPoints(e)
	for i := 0; i < 1000; i++ {
		idx := rand.Int() % max
		me, err := icfp13.MutateAt(e, idx+1, icfp13.GenSexp)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println(me)
		}
	}*/

	/*mutations := icfp13.GenNewRandomUniqGeneration(e)
	  for _,m := range mutations {
	    fmt.Printf("%v --> %X\n",m, icfp13.Eval(m, icfp13.Env{"x_38078":100}))
	  }*/

}
