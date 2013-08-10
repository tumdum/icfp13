package icfp13

import "testing"

func TestScoreOfCorrectProgram(t *testing.T) {
	p := "(lambda (x) (not x))"
	e := Parse([]byte(p))
	cons := GenConstrains(e, 1000)
	score := Score(e, cons)
	if score != 1.0 {
		t.Errorf("expected to have perfect score, but fot '%v'", score)
	}
}
