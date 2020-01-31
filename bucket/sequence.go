package bucket

import (
	"fmt"
	"math"
	"time"
)

type Sequence interface {
	Next() (val int)
}

type Life struct {
	stime    time.Time
	dtime    time.Time
	duration time.Duration
}

func NewLife(d time.Duration, stime, dtime time.Time) Life {
	if !stime.Before(dtime) {
		panic(fmt.Sprintf("failed to init life object, stime must be before the etime stime:%v dtime:%v", stime, dtime))
	}
	return Life{
		stime:    stime,
		dtime:    dtime,
		duration: d,
	}
}

func (l Life) Age() int {
	return int(time.Since(l.stime) / l.duration)
}

func (l Life) IsBorn() bool {
	return time.Now().After(l.stime)
}

func (l Life) IsDead() bool {
	return time.Now().After(l.dtime)
}

// f(x)=c
type Constant struct {
	Life
	Constant int
}

func NewConstant(c int, life Life) Sequence {
	return &Constant{
		Life:     life,
		Constant: c,
	}
}

func (c *Constant) Next() (val int) {
	return c.Constant
}

// f(x)=ax+b
type Linear struct {
	Life
	a float64
	b float64
}

func NewLinear(a, b float64, l Life) Sequence {
	return &Linear{
		Life: l,
		a:    a,
		b:    b,
	}
}

func (l *Linear) Next() (val int) {
	return int(float64(l.Age())*l.a + l.b)
}

// f(x)=x^y+b
type Power struct {
	Life
	y float64
	b float64
}

func NewPower(y, b float64, l Life) Sequence {
	return &Power{
		Life: l,
		y:    y,
		b:    b,
	}
}

func (p *Power) Next() (val int) {
	return int(math.Pow(float64(p.Life.Age()), p.y) + p.b)
}

// f(x)=2^x+b
type Exp2 struct {
	Life
	b float64
}

func NewExp2(b float64, l Life) Sequence {
	return &Exp2{
		Life: l,
		b:    b,
	}
}

func (p *Exp2) Next() (val int) {
	return int(math.Exp2(float64(p.Life.Age())) + p.b)
}
