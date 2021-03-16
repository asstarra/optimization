package lab

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	m "optimization/model"
)

var F, DF func(float64) float64

func Svenn(x, t float64) (float64, float64, int) {
	var f1, f2, f3 float64 = F(x - t), F(x), F(x + t)
	if (f2 >= f1) && (f2 >= f3) {
		return x, x, -1 // bad
	} else if (f2 <= f1) && (f2 <= f3) {
		return x - t, x + t, 0
	} else {
		var del float64
		if f1 > f3 {
			del = t
		} else {
			del = -t
		}
		x1, x2, k := x, x+del, 1
		for ; F(x2) < F(x1); k++ {
			del = del * 2
			x1, x2 = x2, x2+del
		}
		if x < x2 {
			return x, x2, k
		} else {
			return x2, x, k
		}
	}
}

func Half(a, b, eps float64) (float64, int) {
	var c float64
	var iter int = 0
	for c = (a + b) / 2; math.Abs(a-b) > eps; iter++ {
		yk, zk := (a+c)/2, (b+c)/2
		if F(yk) < F(zk) {
			b, c = c, yk
		} else {
			a, c = c, zk
		}
	}
	return c, iter
}

func Gold(a, b, eps float64) (float64, int) {
	gold := (3 - math.Sqrt(5)) / 2
	iter := 0
	for c, d := a+gold*(b-a), b-gold*(b-a); math.Abs(a-b) > eps; iter++ {
		if F(c) < F(d) {
			c, d, b = a+gold*(d-a), c, d
		} else {
			a, c, d = c, d, b-gold*(b-c)
		}
	}
	return (a + b) / 2, iter
}

func Fib(a, b, eps float64, n int) (float64, int) {
	if n < 2 {
		return a, -1
	}
	arrFib := make([]float64, n+1, n+1)
	arrFib[0], arrFib[1] = 1, 1
	for i := 2; i <= n; i++ {
		arrFib[i] = arrFib[i-1] + arrFib[i-2]
	}
	iter := 0
	for c, d := a+(arrFib[n-2]/arrFib[n])*(b-a), a+(arrFib[n-1]/arrFib[n])*(b-a); iter <= n-3 && math.Abs(b-a) > eps; iter++ {
		if F(c) < F(d) {
			c, d, b = a+(arrFib[n-3-iter]/arrFib[n-1-iter])*(d-a), c, d
		} else {
			a, c, d = c, d, c+(arrFib[n-2-iter]/arrFib[n-1-iter])*(b-c)
		}
	}
	return (a + b) / 2, iter
}

func Lab3_1() {
	F = func(a float64) float64 {
		return 10*a*a*a*a + 3*a*a + 2*a - 10*math.Cos(a)
	}
	var t1 time.Time
	rand.Seed(time.Now().UnixNano())
	t1 = time.Now()
	a, b, iter := Svenn(-1, 0.1)
	fmt.Printf("Sven: a=%v, b=%v, iter=%v, dt=%v\n", a, b, iter, time.Since(t1).Nanoseconds())
	h, x1, x2 := 0.001, -100*rand.Float64(), 100*rand.Float64()
	t1 = time.Now()
	a, iter = Half(x1, x2, h)
	fmt.Printf("Half: x=%v, F(x)=%v, iter=%v, dt=%v\n", a, F(a), iter, time.Since(t1).Nanoseconds())
	t1 = time.Now()
	a, iter = Gold(x1, x2, h)
	fmt.Printf("Gold: x=%v, F(x)=%v, iter=%v, dt=%v\n", a, F(a), iter, time.Since(t1).Nanoseconds())
	t1 = time.Now()
	a, iter = Fib(x1, x2, h, 100)
	fmt.Printf("Fib: x=%v, F(x)=%v, iter=%v, dt=%v\n", a, F(a), iter, time.Since(t1).Nanoseconds())

	isMin := false
	t1 = time.Now()
	q, w, iter := m.IntervalSvenn(m.Number(-1), m.Number(0.1), isMin)
	fmt.Printf("Sven: q=%v, w=%v, iter=%v, dt=%v\n", q, w, iter, time.Since(t1).Nanoseconds())
	t1 = time.Now()
	q, iter = m.ExtremumHalf(m.Number(x1), m.Number(x2), h, isMin)
	fmt.Printf("Half: x=%v, F(x)=%v, iter=%v, dt=%v\n", q, q.Value(), iter, time.Since(t1).Nanoseconds())
	t1 = time.Now()
	q, iter = m.ExtremumGold(m.Number(x1), m.Number(x2), h, isMin)
	fmt.Printf("Gold: x=%v, F(x)=%v, iter=%v, dt=%v\n", q, q.Value(), iter, time.Since(t1).Nanoseconds())
	t1 = time.Now()
	q, iter = m.ExtremumFib(m.Number(x1), m.Number(x2), h, isMin)
	fmt.Printf("Fib: x=%v, F(x)=%v, iter=%v, dt=%v\n", q, q.Value(), iter, time.Since(t1).Nanoseconds())
}
