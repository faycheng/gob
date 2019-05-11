package bucket

import "time"

type Bucket interface {
	Get() bool
}

type life struct {
	born     bool
	bornTime time.Time
	duration time.Duration
}

func (l *life) Born() {
	l.born = true
	l.bornTime = time.Now()
}

func (l *life) Age() int {
	return int(time.Since(l.bornTime) / time.Second)
}

func (l *life) IsBorn() bool {
	return l.born
}

func (l *life) IsDead() bool {
	if !l.born {
		return true
	}
	return time.Since(l.bornTime) > l.duration
}

type constantBucket struct {
	*life
	rate int
}

func (b *constantBucket) Get() bool {
	if !b.IsBorn() {
		b.Born()
	}
	if b.IsDead() {
		return false
	}
	time.Sleep(time.Second / time.Duration(b.rate))
	return true
}

type upBucket struct {
	*life
	low  int
	high int
	step int
}

func (b *upBucket) Get() bool {
	if !b.IsBorn() {
		b.Born()
	}
	if b.IsDead() {
		return false
	}
	step := (b.high - b.low) / b.Age()
	curQps := b.low + b.Age()*step
	curSleep := time.Second / time.Duration(curQps)
	time.Sleep(curSleep)
	return true
}

type downBucket struct {
	*life
	low  int
	high int
	step int
}

func (b *downBucket) Get() bool {
	if !b.IsBorn() {
		b.Born()
	}
	if b.IsDead() {
		return false
	}
	step := (b.high - b.low) / b.Age()
	curQps := b.high - b.Age()*step
	curSleep := time.Second / time.Duration(curQps)
	time.Sleep(curSleep)
	return true
}

type rangeBucket struct {
	idx      int
	qpsRange []*constantBucket
}

func (b *rangeBucket) bucket() *constantBucket {
	bucket := b.qpsRange[b.idx]
	if !bucket.IsBorn() {
		bucket.Born()
	}
	if !bucket.IsDead() {
		return bucket
	}
	b.idx++
	if b.idx == len(b.qpsRange) {
		return nil
	}
	bucket = b.qpsRange[b.idx]
	bucket.Born()
	return bucket
}

func (b *rangeBucket) Get() bool {
	bucket := b.qpsRange[b.idx]
	if !bucket.IsBorn() {
		bucket.Born()
	}
	ok := bucket.Get()
	if ok {
		return ok
	}
	b.idx++
	if b.idx == len(b.qpsRange) {
		return false
	}
	bucket = b.qpsRange[b.idx]
	bucket.Born()
	return bucket.Get()
}
