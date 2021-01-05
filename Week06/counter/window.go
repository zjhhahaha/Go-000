package counter

import (
	"time"
)

const defaultSize = 10

type event int

func (e event) oridinal() int {
	return int(e)
}

const (
	Success event = iota
	Failure
	Timeout
	Rejection
)

type bucket struct {
	//TODO:根据event记数，counters[event.oridinal()]
	//counters []int
	counter int
	next    *bucket
}

func (b *bucket) incr() {
	b.counter++
}

func (b *bucket) val() int {
	return b.counter
}

func (b *bucket) reset() {
	b.counter = 0
}

type window struct {
	buckets []bucket
	size    int
	startAt time.Time
}

type windowOpt struct {
	size int
}

func newWindow(opt *windowOpt) *window {
	var size = defaultSize
	if opt != nil {
		size = opt.size
	}
	buckets := make([]bucket, size)
	for offset := range buckets {
		buckets[offset] = bucket{counter: 0}
		nOffset := offset + 1
		if nOffset == size {
			nOffset = 0
		}
		buckets[offset].next = &buckets[nOffset]
	}
	return &window{
		buckets: buckets,
		size:    size,
	}
}

func (w *window) incr(event event, index int) {
	w.buckets[index].incr()
}

func (w *window) resetAll() {
	for offset := range w.buckets {
		w.buckets[offset].reset()
	}
}

func (w *window) reset(offset int) {
	w.buckets[offset].reset()
}

func (w *window) val(offset int) int {
	return w.buckets[offset].val()
}
