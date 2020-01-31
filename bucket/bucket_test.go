package bucket

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucket(t *testing.T) {
	l := NewLife(time.Second, time.Now(), time.Now().Add(time.Second*5))
	c := NewLinear(2, 100, l)
	bucket := NewBucket(l, c)

	for i := 0; i < 3; i++ {
		start := time.Now()
		count := 0
		for time.Since(start) <= time.Second {
			ok := bucket.Get()
			assert.True(t, ok)
			count++
		}
		equal := func(x, y int) bool {
			return int(math.Abs(float64(x-y))) < 10
		}
		assert.True(t, equal(100+2*i, count))
	}
}

func TestBucketLife(t *testing.T) {
	l := NewLife(time.Second, time.Now().Add(time.Second), time.Now().Add(time.Second*3))
	c := NewLinear(2, 100, l)
	bucket := NewBucket(l, c)

	// acquire before the stime
	start := time.Now()
	ok := bucket.Get()
	assert.True(t, ok)
	assert.True(t, time.Since(start) > time.Second && time.Since(start) < time.Second+time.Millisecond*20)

	// acquire after the etime
	time.Sleep(time.Second * 3)
	ok = bucket.Get()
	assert.False(t, ok)
}
