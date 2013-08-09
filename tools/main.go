package main

import (
	"bitbucket.org/tumdum/icfp13"
	"fmt"
)

func main() {
	fmt.Println(icfp13.Parse([]byte("(lambda (x) (fold x 0 (lambda (y z) (or y z))))")))
	fmt.Println(icfp13.Parse([]byte("(lambda (x_2131) (fold (shr1 (not x_2131)) x_2131 (lambda (x_2132 x_2133) (shl1 (xor x_2132 x_2133)))))")))
	fmt.Println(icfp13.Parse([]byte("(lambda (x_10331) (fold x_10331 0 (lambda (x_10331 x_10332) (plus (shr1 (xor 0 x_10331)) x_10332))))")))
	// s := icfp13.Parse([]byte("(lambda (x) (fold x 0 (lambda (y z) (or y z))))"))
	simple := icfp13.Parse([]byte("(or (not x) x)"))
	// fmt.Println(icfp13.Eval(s, icfp13.Env{"x":0x1122334455667788}))
	fmt.Printf("simple: %X\n", icfp13.Eval(simple, icfp13.Env{"x": 0xFFFFFFFFFFFFFFFF}))
	f := icfp13.Parse([]byte("(or a (or b (or c (or d (or e (or f (or g (or h 0))))))))"))
	fmt.Printf("f: %X\n", icfp13.Eval(f, icfp13.Env{"a": 0x11, "b": 0x22, "c": 0x33, "d": 0x44, "e": 0x55, "f": 0x66, "g": 0x77, "h": 0x88}))

	e := icfp13.Parse([]byte("(lambda (x) x)"))
	fmt.Println(e, icfp13.Eval(e, icfp13.Env{"x": 42}))
}
