package bucket

import (
	"time"
)

type Bucket struct {
	life             Life
	rate             int
	sleep            time.Duration
	lastGenerateTime time.Time
	lastGetTime      time.Time
	sequence         Sequence
}

func NewBucket(l Life, seq Sequence) *Bucket {
	return &Bucket{
		life:     l,
		sequence: seq,
	}
}

func (b *Bucket) Rate() (rate int, ok bool) {
	for !b.life.IsBorn() {
		time.Sleep(time.Second)
	}
	if b.life.IsDead() {
		return 0, false
	}
	if b.lastGenerateTime.IsZero() || time.Since(b.lastGenerateTime) >= time.Second {
		b.rate = b.sequence.Next()
		b.lastGenerateTime = time.Now()
	}
	for b.rate == 0 {
		if b.life.IsDead() {
			return 0, false
		}
		time.Sleep(time.Second)
		b.rate = b.sequence.Next()
		b.lastGenerateTime = time.Now()
	}
	return b.rate, true
}

func (b *Bucket) Get() bool {
	rate, ok := b.Rate()
	if !ok {
		return false
	}
	now := time.Now()
	if b.lastGetTime.IsZero() {
		b.lastGetTime = now
	}
	sleepSpan := time.Second / time.Duration(rate)
	// sleep represents how much time we should sleep.
	// Inspired by github.com/user-go/ratelimit.
	b.sleep += sleepSpan - time.Since(b.lastGetTime)
	if b.sleep > 0 {
		// span: system call time of time.Sleep and time.Add
		// other: time of caller using
		// span + other = time.Now().Sub(b.last)
		time.Sleep(b.sleep)
		b.lastGetTime = now.Add(b.sleep)
		b.sleep = 0
		return true
	}
	b.lastGetTime = now
	return true
}
