package counter

import (
	"sync"
	"time"
)

type rollingCounter struct {
	mu       sync.Mutex
	window   *window
	duration time.Duration
	size     int
	offset   int
	lastIncr time.Time
	nowFunc  func() time.Time
}

func NewRollingCounter(size int, duration time.Duration) *rollingCounter {
	return &rollingCounter{
		mu: sync.Mutex{},
		window: newWindow(&windowOpt{
			size: size,
		}),
		size:     size,
		duration: duration,
		offset:   0,
		nowFunc:  time.Now,
		lastIncr: time.Now(),
	}
}

func (c *rollingCounter) Incr(event event) {
	c.mu.Lock()
	defer c.mu.Unlock()
	incrOffset := c.getIncrOffeset()
	currentOffset := incrOffset + c.offset
	//超出长度，全部重置，放入第一个bukcet，更新时间
	//等于0，则直接插入
	//大于0，需要将无效的重置，插入数据，并且更新时间
	if incrOffset > c.size {
		c.window.resetAll()
		c.offset = 0
	} else {
		c.offset = currentOffset % c.size
		for i := 1; i <= incrOffset; i++ {
			offset := c.offset - i
			if offset < 0 {
				offset = c.size - offset
			}
			c.window.reset(offset)
		}
		c.offset = currentOffset
	}
	c.window.incr(event, c.offset)
	c.lastIncr = c.lastIncr.Add(time.Duration(incrOffset) * c.duration)
}

func (c *rollingCounter) getIncrOffeset() int {
	now := c.nowFunc()
	return int(now.Sub(c.lastIncr) / c.duration)
}

func (c *rollingCounter) GetCurrentCount(event event) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.window.val(c.offset)

}
