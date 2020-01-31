package bucket

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstant(t *testing.T) {
	l := NewLife(time.Second, time.Now(), time.Now().Add(time.Second*10))
	c := NewConstant(100, l)

	assert.Equal(t, 100, c.Next())
}

func TestLiner(t *testing.T) {
	l := NewLife(time.Second, time.Now(), time.Now().Add(time.Second*3))
	c := NewLinear(2, 100, l)

	for i := 0; i < 3; i++ {
		assert.Equal(t, 100+2*i, c.Next())
		time.Sleep(time.Second)
	}
}

func TestPower(t *testing.T) {
	l := NewLife(time.Second, time.Now(), time.Now().Add(time.Second*3))
	p := NewPower(2, 100, l)

	for i := 0; i < 3; i++ {
		assert.Equal(t, int(100+math.Pow(float64(i), 2)), p.Next())
		time.Sleep(time.Second)
	}
}

func TestExp2(t *testing.T) {
	l := NewLife(time.Second, time.Now(), time.Now().Add(time.Second*3))
	p := NewExp2(100, l)

	for i := 0; i < 3; i++ {
		assert.Equal(t, int(100+math.Exp2(float64(i))), p.Next())
		time.Sleep(time.Second)
	}
}
