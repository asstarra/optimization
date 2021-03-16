package model

import (
	"fmt"
	"math"
	"time"
)

type Function interface {
	Value() float64
	Sum(b Function) Function
	Sub(b Function) Function
	Mul(f float64) Function
	Minus() Function
	Dist(b Function) float64
	IsDistZero(b Function, eps float64) bool
	Len() int
	GetCoordinate(i int) float64
	SetCoordinate(i int, val float64) Function
	Zero() Function
	E() Function
	String() string
	Copy() Function
	Less(b Function, isMin bool) bool
	// Equal(b Function) bool
}

func Extremum(arr []Function, isMin bool) Function {
	a := arr[0]
	for _, val := range arr {
		if (val.Value() < a.Value()) == isMin {
			a = val
		}
	}
	return a.Copy()
}

func Equal(a, b Function) bool {
	flag := a.Len() == b.Len()
	for i := 0; flag && i < a.Len(); i++ {
		flag = a.GetCoordinate(i) == b.GetCoordinate(i)
	}
	return flag
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	fmt.Printf("%v: %v\n", msg, time.Since(start))
}

// ------------------- Number ------------------- //

type Number float64

func (x Number) Value() float64 {
	a := float64(x)
	return -(10*a*a*a*a + 3*a*a + 2*a - 10*math.Cos(a))
}

func (a Number) Sum(b Function) Function {
	b2, ok := b.(Number)
	if !ok {
		panic("not_correct")
	}
	return Number(float64(a) + float64(b2))
}

func (a Number) Sub(b Function) Function {
	b2, ok := b.(Number)
	if !ok {
		panic("not_correct")
	}
	return Number(float64(a) - float64(b2))
}

func (a Number) Mul(f float64) Function {
	return Number(float64(a) * f)
}

func (a Number) Minus() Function {
	return Number(-float64(a))
}

func (a Number) Dist(b Function) float64 {
	b2, ok := b.(Number)
	if !ok {
		panic("not_correct")
	}
	return float64(a) - float64(b2)
}

func (a Number) IsDistZero(b Function, eps float64) bool {
	return math.Abs(a.Dist(b)) <= eps
}

func (a Number) Len() int {
	return 1
}

func (a Number) GetCoordinate(i int) float64 {
	if i < 0 || i >= a.Len() {
		panic("out of size array")
	}
	return float64(a)
}

func (a Number) SetCoordinate(i int, val float64) Function {
	if i < 0 || i >= a.Len() {
		panic("out of size array")
	}
	return Number(val)
}

func (a Number) Zero() Function {
	return Number(0)
}

func (a Number) E() Function {
	return Number(1)
}

func (a Number) String() string {
	return fmt.Sprint(float64(a))
}

func (a Number) Copy() Function {
	return Number(float64(a))
}

func (a Number) Less(b Function, isMin bool) bool {
	return (a.Value() < b.Value() && isMin) || (a.Value() > b.Value() && !isMin)
}

// ------------------- Vector ------------------- //

type Vector struct {
	xi []float64
}

var FV func(v Vector) float64

func NewVector(arr []float64) Vector {
	return Vector{xi: arr}
}

func (v Vector) Value() float64 {
	return FV(v)
}
func (v Vector) Sum(b Function) Function {
	if v.Len() != b.Len() {
		panic("arrays of different lenghts")
	}
	w := Vector{xi: make([]float64, v.Len(), v.Len())}
	for i, a := range v.xi {
		w.xi[i] = a + b.GetCoordinate(i)
	}
	return w
}
func (v Vector) Sub(b Function) Function {
	return v.Sum(b.Minus())
}
func (v Vector) Mul(f float64) Function {
	w := Vector{xi: make([]float64, v.Len(), v.Len())}
	for i, a := range v.xi {
		w.xi[i] = a * f
	}
	return w
}
func (v Vector) Minus() Function {
	w := Vector{xi: make([]float64, v.Len(), v.Len())}
	for i, a := range v.xi {
		w.xi[i] = -a
	}
	return w
}
func (v Vector) Dist(b Function) float64 {
	if v.Len() != b.Len() {
		panic("arrays of different lenghts")
	}
	var dist float64 = 0
	for i, val := range v.xi {
		dist += math.Pow(val-b.GetCoordinate(i), 2)
	}
	return math.Sqrt(dist)
}
func (v Vector) IsDistZero(b Function, eps float64) bool {
	// fmt.Println("-7-", v, b, v.Dist(b), eps)
	return v.Dist(b) <= eps
}
func (v Vector) Len() int {
	return len(v.xi)
}
func (v Vector) GetCoordinate(i int) float64 {
	if i < 0 || i >= v.Len() {
		panic("out of size array")
	}
	return v.xi[i]
}

func (v Vector) SetCoordinate(i int, val float64) Function {
	if i < 0 || i >= v.Len() {
		panic("out of size array")
	}
	w := Vector{xi: make([]float64, v.Len(), v.Len())}
	for i, f := range v.xi {
		w.xi[i] = f
	}
	w.xi[i] = val
	return w
}

func (v Vector) Zero() Function {
	return Vector{xi: make([]float64, v.Len(), v.Len())}
}

func (v Vector) E() Function {
	w := Vector{xi: make([]float64, v.Len(), v.Len())}
	for i, _ := range w.xi {
		w.xi[i] = 1
	}
	return w
}

func (v Vector) String() string {
	return fmt.Sprintf("{pos=%v, f=%f}", v.xi, v.Value())
}

func (v Vector) Copy() Function {
	w := Vector{xi: make([]float64, v.Len(), v.Len())}
	copy(w.xi, v.xi)
	return w
}

func (v Vector) Less(b Function, isMin bool) bool {
	return (v.Value() < b.Value() && isMin) || (v.Value() > b.Value() && !isMin)
}
