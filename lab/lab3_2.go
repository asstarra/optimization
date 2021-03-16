package lab

import (
	"fmt"
	"math"
)

func Extremum(arr []float64, isMin bool) float64 {
	if len(arr) == 0 {
		panic("array is null")
	}
	a := arr[0]
	for _, val := range arr {
		if (F(a) > F(val)) == isMin {
			a = val
		}
	}
	return a
}

func MinMax(arr []float64, isMin bool) float64 {
	if len(arr) == 0 {
		panic("array is null")
	}
	a := arr[0]
	for _, val := range arr {
		if (a > val) == isMin {
			a = val
		}
	}
	return a
}

func Belong(x1, x2, x3 float64) bool {
	return (x1 <= x3 && x3 <= x2) || (x1 >= x3 && x3 >= x2)
}

func QuadExtremumPov(x1, x2, x3, pov float64) float64 {
	return F(x1)*(math.Pow(x2, pov)-math.Pow(x3, pov)) +
		F(x2)*(math.Pow(x3, pov)-math.Pow(x1, pov)) +
		F(x3)*(math.Pow(x1, pov)-math.Pow(x2, pov))
}

func QuadInter(x, t, eps, d float64, isMin bool, maxIter int) (float64, int) {
	defer duration(track("QuadInter"))
	nextX := func(x1 float64) (float64, float64, float64) {
		x2 := x1 + t
		if (F(x2) < F(x1)) == isMin {
			return x1, x2, x2 + 2*t
		} else {
			return x1, x2, x1 - 2*t
		}
	}
	iter := 0
	for x1, x2, x3 := nextX(x); ; iter++ {
		xe := Extremum([]float64{x1, x2, x3}, isMin)
		chis, znam := QuadExtremumPov(x1, x2, x3, 2), QuadExtremumPov(x1, x2, x3, 1)
		if znam == 0 {
			if x1 == xe {
				panic("error alg")
			}
			x1, x2, x3 = nextX(x1)
		} else {
			xs := chis / (znam * 2)
			feps, xeps := math.Abs((F(xe)-F(xs))/F(xs)), math.Abs((xe-xs)/xs)
			if (feps <= eps && xeps <= d) || iter >= maxIter {
				return xs, iter
			}
			if Belong(x1, x3, xs) {
				xs = Extremum([]float64{xs, xe}, isMin)
			}
			x1, x2, x3 = nextX(xs)
		}
	}
}

func QuadInterMod(x, t, eps, d float64, isMin bool, maxIter int) (float64, int) {
	defer duration(track("QuadInterMod"))
	nextX := func(x1 float64) (float64, float64, float64, float64) {
		x2 := x1 + t
		if (F(x2) < F(x1)) == isMin {
			return x1, x2, x2 + 2*t, 2 * t
		} else {
			return x1, x2, x1 - 2*t, -2 * t
		}
	}
	iter := 0
	for x1, x2, x3, del := nextX(x); ; iter++ {
		xe := Extremum([]float64{x1, x2, x3}, isMin)
		chis, znam := QuadExtremumPov(x1, x2, x3, 2), QuadExtremumPov(x1, x2, x3, 1)
		if znam == 0 {
			if x1 == xe {
				panic("error alg")
			}
			x1, x2, x3, del = nextX(x1)
		} else {
			xs := chis / (znam * 2)
			if (del < 0 && xs > x2) || (del > 0 && xs < x1) {
				x1 = MinMax([]float64{x1, x2, x3}, del < 0)
				// x2, x3 = x1+del, x1+2*del
				// for i := 0; (F(x1)-F(x2)) <= (F(x2)-F(x3)) == isMin && i < maxIter; i, iter = i+1, iter+1 {
				// 	x1, x2, x3 = x2, x3, x3+del
				// }
				x2, x3 = x1+del, x1+3*del
				for i := 0; (F(x1)-F(x2)) <= (F(x2)-F(x3))/2 == isMin && i < maxIter; i, iter = i+1, iter+1 {
					x1, x2, x3, del = x2, x3, x3+4*del, 2*del
				}
				x1, x2, x3, del = nextX(x3)
			} else {
				feps, xeps := math.Abs((F(xe)-F(xs))/F(xs)), math.Abs((xe-xs)/xs)
				if (feps <= eps && xeps <= d) || iter >= maxIter {
					return xs, iter
				}
				if Belong(x1, x3, xs) {
					xs = Extremum([]float64{xs, xe}, isMin)
				}
				x1, x2, x3, del = nextX(xs)
			}
		}
	}
}

func Sign(x float64) float64 {
	if x == 0 {
		return 0
	} else if x > 0 {
		return 1
	} else {
		return -1
	}
}

func CubInter(x, t, eps, d float64, isMin bool, maxIter int) (float64, int) {
	defer duration(track("CubInter"))
	var signE float64 = -1
	if isMin {
		signE = 1
	}
	t = -1 * signE * Sign(DF(x)) * t
	var x1, x2 float64 = x, x + t
	for k := 1; DF(x1)*DF(x2) > 0; k++ {
		x1, x2 = x2, x2+math.Pow(2, float64(k))*t
	}
	var f1, f2, df1, df2 float64 = signE * F(x1), signE * F(x2), signE * DF(x1), signE * DF(x2)
	for iter := 1; ; iter++ {
		z := (3*(f1-f2)/(x2-x1) + df1 + df2)
		w := math.Pow(z*z-df1*df2, 0.5)
		if x1 > x2 {
			w = -w
		}
		m := (df2 + w - z) / (df2 - df1 + 2*w)
		var xs, fxs, dfxs float64
		if m < 0 {
			xs = x2
		} else if m > 1 {
			xs = x1
		} else {
			xs = x2 - m*(x2-x1)
		}
		for fxs = signE * F(xs); fxs > f1; {
			xs = xs - 0.5*(xs-x1)
			fxs = signE * F(xs)
		}
		dfxs = signE * DF(xs)

		if (dfxs <= eps && math.Abs((xs-x1)/xs) < d) || iter >= maxIter {
			return xs, iter
		} else if dfxs*df1 < 0 {
			x1, f1, df1, x2, f2, df2 = xs, fxs, dfxs, x1, f1, df1
		} else {
			x1, f1, df1, x2, f2, df2 = x2, f2, df2, xs, fxs, dfxs
		}
	}
}

func Lab3_2() {
	fmt.Println()
	maxIter, x1, step, eps, d := 100, -1.0, 0.01, 0.001, 0.001
	F = func(a float64) float64 {
		return a * (a + 2) * (a - 2)
	}
	DF = func(x float64) float64 {
		x1, x2 := x-0.01, x+0.01
		return (F(x2) - F(x1)) / (x2 - x1)
	}

	x, iter := QuadInter(x1, step, eps, d, true, maxIter)
	fmt.Println(x, F(x), iter, "\n")
	x, iter = QuadInterMod(x1, step, eps, d, true, maxIter)
	fmt.Println(x, F(x), iter, "\n")

	// F = func(a float64) float64 {
	// 	return -(10*a*a*a*a + 3*a*a + 2*a - 10*math.Cos(a))
	// }
	x, iter = QuadInter(x1, step, eps, d, false, maxIter)
	fmt.Println(x, F(x), iter, "\n")
	x, iter = CubInter(x1, step, eps, d, false, maxIter)
	fmt.Println(x, F(x), iter, "\n")
	// F = func(a float64) float64 {
	// 	return (10*a*a*a*a + 3*a*a + 2*a - 10*math.Cos(a))
	// }
	x, iter = CubInter(x1, step, eps, d, true, maxIter)
	fmt.Println(x, F(x), iter, "\n")

	F = func(a float64) float64 {
		return 10*a*a*a*a + 3*a*a + 2*a - 10*math.Cos(a)
	}
	x, iter = QuadInter(x1, step, eps, d, true, maxIter)
	fmt.Println(x, F(x), iter, "\n")
	x, iter = CubInter(x1, step, eps, d, true, maxIter)
	fmt.Println(x, F(x), iter, "\n")
}
