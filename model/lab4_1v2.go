package model

import (
	// "math"
	"fmt"
)

// func ExtremumLine(x0, t Function, eps float64, isMin bool, f func(a, b Function, eps float64, isMin bool) (Function, int)) (Function, int) {

// }

func Hook_Jeeves2(x0, t Function, eps, betta float64, isMin bool) (Function, int) {
	for iter := 0; ; iter++ {
		var dx Function = t.Zero() // Шаг 2
		for i := 0; i < t.Len(); i++ {
			dxi := t.Zero().SetCoordinate(i, t.GetCoordinate(i))
			if (x0.Value() > x0.Sum(dxi).Value()) == isMin {
				dx = dx.Sum(dxi)
			} else if (x0.Value() > x0.Sub(dxi).Value()) == isMin {
				dx = dx.Sub(dxi)
			}
		}
		var i int
		a, b, i := IntervalSvenn(x0, dx, isMin)
		if i == -1 {
			t = t.Mul(1 / betta)
			fmt.Println("1 Warning")
			continue
		}
		a, i = ExtremumGold(a, b, dx.Zero().Dist(dx), isMin)
		if a.IsDistZero(x0, eps) {
			return a, iter
		}
		if t.Zero().IsDistZero(t, eps) {
			t = t.Mul(1 / betta)
		}
		x0 = a
	}
}

/*type Symplex struct {
	points []Function
}

func NewSymplex(x0 Function, step float64) Symplex { // шаг 2
	var n = x0.Len() + 1
	points := make([]Function, n, n)
	points[0] = x0
	nf := float64(n)
	l2 := x0.E().Mul(step * (math.Sqrt(nf+1) - 1) / nf / math.Sqrt2)
	for i := 1; i < n; i++ {
		points[i] = l2.SetCoordinate(i, step*(math.Sqrt(nf+1)+nf-1)/nf/math.Sqrt2).Sum(x0)
	}
	return Symplex{points: points}
}

func (s Symplex) Len() int {
	return len(s.points)
}

func (s Symplex) Extremum(isMin bool) Function {
	return MinMax(s.points, isMin)
}

func (s Symplex) Centre(xh Function) Function {
	xc := xh.Minus()
	for _, val := range s.points {
		xc = xc.Sum(val)
	}
	return xc
}

// func Nelder_Mid(x0 Function, step float64, isMin bool) {
// 	s := NewSymplex(x0, step)                       // шаг 2
// 	xl, xh := s.Extremum(isMin), s.Extremum(!isMin) // шаг 3
// 	xc := s.Centre(xh)
// }
*/
