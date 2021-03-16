package lab

import (
	"fmt"
	"math"
	"time"

	m "optimization/model"
)

func Lab4_1() {
	m.FV = func(v m.Vector) float64 {
		n := v.Len()
		var f, a, b float64 = 100, 150, 2
		for i := 0; i < n-1; i++ {
			f += a*math.Pow(math.Pow(v.GetCoordinate(i), 2)-v.GetCoordinate(i+1), 2) + b*math.Pow(v.GetCoordinate(i)-1, 2)
		}
		return f
	}
	var step, epsf, epsx float64 = 1, 0.000001, 0.00001
	var mir, mul, dif, red, twist float64 = 1, 2, 0.5, 0.5, 0.01
	var isMin, maxIter = true, 10000
	t := m.NewVector([]float64{step, step, step})
	// t := m.NewVector([]float64{step, step})
	x0 := t.Mul(-10)
	// x0 := m.NewVector([]float64{-100, 150, 0})
	t1 := time.Now()
	xe, iter := m.Hook_Jeeves(x0, t, epsf, mul, isMin)
	fmt.Println(xe, xe.Value(), iter, time.Since(t1))
	// t1 = time.Now()
	// xe, iter = m.Nelder_Mid(x0, epsf, epsx, step, 1/betta, twist, isMin, maxIter)
	// fmt.Println(xe, xe.Value(), iter, time.Since(t1))
	// xe, iter = m.Hook_Jeeves2(x0, t, epsf, betta, true)
	// fmt.Println(xe, xe.Value(), iter)
	t1 = time.Now()
	xe, iter = m.Nelder_Mid2(x0, epsf, epsx, step, mir, mul, dif, red, twist, isMin, maxIter)
	fmt.Println(xe, xe.Value(), iter, time.Since(t1))
}
