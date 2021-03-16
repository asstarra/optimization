package model

import (
	"fmt"
	"math"
)

func ExtremumLine(x0, dx Function, eps float64, isMin bool, f func(a, b Function, eps float64, isMin bool) (Function, int)) (Function, int) {
	a, b, iter := IntervalSvenn(x0, dx, isMin)
	if iter == -1 {
		fmt.Println(a, b, iter)
		a, b, iter = IntervalSvenn(b, dx, isMin)
		fmt.Println(a, b, iter)
		if iter == -1 {
			return x0, -1
		}
	}
	return f(a, b, eps, isMin)
}

func Hook_Jeeves(x0, t Function, epsx, betta float64, isMin bool) (Function, int) {
	var x = x0.Copy()
	for iter := 0; ; iter++ {
		var dx Function = t.Zero() // Шаг 2
		for i := 0; i < t.Len(); i++ {
			dxi := t.Zero().SetCoordinate(i, t.GetCoordinate(i))
			if (x.Value() > x.Sum(dxi).Value()) == isMin {
				dx = dx.Sum(dxi)
			} else if (x.Value() > x.Sub(dxi).Value()) == isMin {
				dx = dx.Sub(dxi)
			}
		}
		if dx.Zero().IsDistZero(dx, 0) {
			t = t.Mul(1 / betta)
			continue
		}

		a, i := ExtremumLine(x, dx, dx.Zero().Dist(dx), isMin, ExtremumGold)
		if i == -1 {
			t = t.Mul(1 / betta)
			fmt.Println("2 Warning")
			iter++
			continue
		}

		if a.IsDistZero(x, epsx) {
			return a, iter
		}
		// if t.Zero().IsDistZero(t, epsx) {
		// 	t = t.Mul(1 / betta)
		// }
		x = a
	}
}

type Symplex struct {
	points []Function
}

// Построение нового симплекса с фиксированным шагом.
func NewSymplex(x0 Function, step float64) Symplex { // шаг 2.
	var n = x0.Len() + 1
	points := make([]Function, n, n)
	points[0] = x0.Copy()
	nf := float64(n)
	l2 := x0.E().Mul(step * (math.Sqrt(nf+1) - 1) / nf / math.Sqrt2)
	for i := 1; i < n; i++ {
		points[i] = l2.SetCoordinate(i-1, step*(math.Sqrt(nf+1)+nf-1)/nf/math.Sqrt2).Sum(x0)
	}
	return Symplex{points: points}
}

// Размерность симплекса.
func (s Symplex) Len() int {
	return len(s.points)
}

// Возврашает точку симлекса в которой функция принмает
// минимальное (максимальное) значение.
func (s Symplex) Extremum(isMin bool) Function { // шаг 3.
	return Extremum(s.points, isMin)
}

// Ищем центр симплекса с исключением точки xh.
func (s Symplex) Centre(xh Function) Function { // шаг 4.
	xc := xh.Minus()
	for _, val := range s.points {
		xc = xc.Sum(val)
	}
	return xc.Mul(1 / float64(s.Len()-1))
}

func (s Symplex) GetIndex(x Function) int {
	for i, val := range s.points {
		if Equal(val, x) {
			return i
		}
	}
	return -1
}

func (s Symplex) Reduction(k float64, isMin bool) Symplex { // шаг 7.
	xl := s.Extremum(isMin)
	for i, val := range s.points {
		s.points[i] = xl.Sum(val.Sub(xl).Mul(k))
	}
	return s
}

func (s Symplex) Replacement(xh, xn Function) Symplex { // шаг 6.
	for i, val := range s.points {
		if Equal(val, xh) {
			s.points[i] = xn.Copy()
		}
	}
	return s
}

func (s Symplex) IsExtremum(eps float64, isMin bool) bool { // шаг 8.
	xl := s.Extremum(isMin)
	var sum float64 = 0
	for i := 0; i < s.Len(); i++ {
		sum += math.Pow(xl.Value()-s.points[i].Value(), 2)
	}
	sum = math.Sqrt(sum / float64(s.Len()))
	return sum < eps
}

func (s Symplex) IsTwisted(k float64) bool {
	arccos := func(a, b, c Function) float64 {
		x, y := a.Sub(b), a.Sub(c)
		var sum float64 = 0
		for i := 0; i < a.Len(); i++ {
			sum += x.GetCoordinate(i) * y.GetCoordinate(i)
		}
		return math.Acos(math.Abs(sum / x.Zero().Dist(x) / y.Zero().Dist(y)))
	}
	var min float64 = 180
	for i := 0; i < s.Len(); i++ {
		for j := 0; j < s.Len(); j++ {
			for k := j + 1; k < s.Len(); k++ {
				if i != j && i != k {
					acos := arccos(s.points[i], s.points[j], s.points[k])
					if acos < min {
						min = acos
					}
				}
			}
		}
	}
	return min < k
}

func (s Symplex) String() string {
	str := "{"
	for i := 0; i < s.Len(); i++ {
		str += "\n" + s.points[i].String()
	}
	return str + "}\n"
}

// func Nelder_Mid(x0 Function, eps, d, stepn, reduc, twist float64, isMin bool, maxIter int) (Function, int) {
// 	step := stepn
// 	s := NewSymplex(x0, step) // шаг 2.
// 	for iter := 0; iter < maxIter; iter++ {
// 		if step < d/10 {
// 			step = stepn
// 		}
// 		xl, xh := s.Extremum(isMin), s.Extremum(!isMin) // шаг 3.
// 		if xl.IsDistZero(xh, d) {
// 			return xl, iter
// 		}
// 		xc := s.Centre(xh) // шаг 4.
// 		// fmt.Println("-6-", xc.Sub(xh).Mul(0.1).Dist(xc.Zero()))
// 		xn, i := ExtremumLine(xc, xc.Sub(xh), d, isMin, ExtremumGold)                // шаг 5.
// 		if i == -1 || xn.Value() == xh.Value() || xn.Value() > xh.Value() == isMin { // xn - хуже xh.
// 			// fmt.Printf("-1- xn=%v, f(xn)=%v, f(xl)=%v, xl=%v\n", xn, xn.Value(), xl.Value(), xl)
// 			step = step / 2
// 			s = s.Reduction(reduc, isMin) // шаг 7.
// 		} else { // xn - лучше xh.
// 			// fmt.Printf("-2- xn=%v, f(xn)=%v, f(xl)=%v, xl=%v\n", xn, xn.Value(), xl.Value(), xl)
// 			s = s.Replacement(xh, xn)                   // шаг 6.
// 			if xh = s.Extremum(!isMin); Equal(xh, xn) { // xn - худшая точка.
// 				step = step / 2
// 				// fmt.Printf("-4- s=%v", s)
// 				s = s.Reduction(reduc, isMin) // шаг 7.
// 			}
// 			// fmt.Println(i, xn.IsDistZero(xl, d), s.IsExtremum(eps, isMin))
// 			if xn.IsDistZero(xl, d) && s.IsExtremum(eps, isMin) { // шаг 8.
// 				return xn, iter
// 			}
// 			// xl = xn
// 		}

// 		if /*iter%10 == 9 && */ s.IsTwisted(twist) { // шаг9.
// 			// fmt.Println("-3-")
// 			xl := s.Extremum(isMin)
// 			s = NewSymplex(xl, step)
// 		}
// 		// fmt.Printf("iter=%d, xl=%v, f(xl)=%f\n, s=%v\n\n", iter, xl, xl.Value(), s)
// 	}
// 	return s.Extremum(isMin), maxIter
// }

func Nelder_Mid2(x0 Function, epsf, epsx, step, mir, mul, dif, red, twist float64, isMin bool, maxIter int) (Function, int) {
	if epsf <= 0 || epsx <= 0 || step <= 0 || mir < 1 || mul <= 1 ||
		dif >= 1 || dif <= 0 || red >= 1 || red <= 0 || twist <= 0 {
		return x0, -1
	}
	s := NewSymplex(x0, step) // шаг 2.
	for iter := 0; iter < maxIter; iter++ {
		if /*iter%10 == 9 &&*/ s.IsTwisted(twist) { // шаг9.
			// fmt.Println("-9- ti\n\n")
			xl := s.Extremum(isMin)
			s = NewSymplex(xl, step)
		}
		xl, xh := s.Extremum(isMin), s.Extremum(!isMin) // шаг 3.
		xc := s.Centre(xh)                              // шаг 4.
		xn := xc.Sub(xh).Mul(mir).Sum(xc)               // шаг 5.
		// fmt.Println(xc, xc.Sub(xh), xc.Sub(xh).Mul(mir), isMin)
		if (xn.Value() < xl.Value()) == isMin { // шаг 6.
			// fmt.Printf("-1- xn=%v, xc=%v\n", xn, xc)
			for xk := xc.Sub(xh).Mul(mul).Sum(xn); xk.Value() < xn.Value() == isMin; xk = xc.Sub(xh).Mul(mul).Sum(xk) {
				xn = xk
			}
			// fmt.Printf("-2- xn=%v, xh=%v, xl=%v\n", xn, xh, xl)
			xl = xn
			s = s.Replacement(xh, xn)
		} else if (xn.Value() > xh.Value()) == isMin {
			xn = xc.Sub(xh).Mul(dif).Sum(xc)
			// fmt.Printf("-3- xn=%v, xc=%v, xh=%v\n", xn, xc, xh)
			if (xn.Value() > xh.Value()) == isMin {
				// fmt.Println("-4-")
				s = s.Reduction(red, isMin) // шаг 7.
			} else {
				// fmt.Println("-5-")
				s = s.Replacement(xh, xn)
			}
			xl = s.Extremum(isMin)
		} else {
			// fmt.Printf("-6- xn=%v, xh=%v, xl=%v\n", xn, xh, xl)
			s = s.Replacement(xh, xn)
		}
		// fmt.Printf("iter=%d, xl=%v, f(xl)=%f\n, s=%v\n\n", iter, xl, xl.Value(), s)
		if xl.IsDistZero(xh, epsx) && s.IsExtremum(epsf, isMin) { // шаг 8.
			return xn, iter
		}
	}
	return s.Extremum(isMin), maxIter
}
