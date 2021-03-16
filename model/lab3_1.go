package model

import (
	"fmt"
	"math"
)

func ExtremumHalf(a, b Function, epsx float64, isMin bool) (Function, int) {
	half := func(a, b Function) Function {
		return a.Sum(b).Mul(0.5)
	}
	var c Function
	var iter int = 0
	for c = half(a, b); !a.IsDistZero(b, epsx); iter++ {
		yk, zk := half(a, c), half(b, c)
		if (c.Value() < yk.Value()) == isMin && (c.Value() < zk.Value()) == isMin {
			a, b = yk, zk
		} else if (yk.Value() < zk.Value()) == isMin {
			b, c = c, yk
		} else {
			a, c = c, zk
		}
	}
	return c, iter
}

func ExtremumGold(a, b Function, epsx float64, isMin bool) (Function, int) {
	var gold float64 = (3 - math.Sqrt(5)) / 2
	var iter int = 0
	for c, d := a.Sum(b.Sub(a).Mul(gold)), b.Sub(b.Sub(a).Mul(gold)); !a.IsDistZero(b, epsx); iter++ {
		if (c.Value() < d.Value()) == isMin {
			c, d, b = a.Sum(d.Sub(a).Mul(gold)), c, d
		} else {
			a, c, d = c, d, b.Sub(b.Sub(c).Mul(gold))
		}
	}
	return a.Sum(b).Mul(0.5), iter
}

func ExtremumFib(a, b Function, epsx float64, isMin bool) (Function, int) {
	// var n int = int(math.Abs(a.Dist(b)) / epsx)
	var n int = int(1 / epsx)
	fmt.Println(a.Dist(b), epsx, n)
	if n < 2 {
		return a, -1
	}
	arrFib := make([]float64, n+1, n+1)
	arrFib[0], arrFib[1] = 1, 1
	for i := 2; i <= n; i++ {
		arrFib[i] = arrFib[i-1] + arrFib[i-2]
	}
	iter := 0
	for c, d := a.Sum(b.Sub(a).Mul(arrFib[n-2]/arrFib[n])), a.Sum(b.Sub(a).Mul(arrFib[n-1]/arrFib[n])); iter <= n-3 && !a.IsDistZero(b, epsx); iter++ {
		if (c.Value() < d.Value()) == isMin {
			c, d, b = a.Sum(d.Sub(a).Mul(arrFib[n-3-iter]/arrFib[n-1-iter])), c, d
		} else {
			a, c, d = c, d, c.Sum(b.Sub(c).Mul(arrFib[n-2-iter]/arrFib[n-1-iter]))
		}
	}
	return a.Sum(b).Mul(0.5), iter
}

func IntervalSvenn(x, t Function, isMin bool) (Function, Function, int) {
	var f1, f2, f3 float64 = x.Sub(t).Value(), x.Value(), x.Sum(t).Value()
	if (isMin && (f2 >= f1) && (f2 >= f3)) || (!isMin && (f2 <= f1) && (f2 <= f3)) {
		return x, Extremum([]Function{x.Sub(t), x, x.Sum(t)}, isMin), -1 // bad
	} else if (isMin && (f2 <= f1) && (f2 <= f3)) || (!isMin && (f2 >= f1) && (f2 >= f3)) {
		return x.Sub(t), x.Sum(t), 0
	} else {
		var del Function
		if (isMin && f1 > f3) || (!isMin && f1 < f3) {
			del = t.Copy()
		} else {
			del = t.Minus()
		}
		x1, x2, k := x, x.Sum(del), 1
		for ; (isMin && x2.Value() < x1.Value()) || (!isMin && x2.Value() > x1.Value()); k++ {
			del = del.Mul(2)
			x1, x2 = x2, x2.Sum(del)
		}
		// if x < x2 {
		// 	return x, x2, k
		// } else {
		// 	return x2, x, k
		// }
		return x, x2, k
	}
}
