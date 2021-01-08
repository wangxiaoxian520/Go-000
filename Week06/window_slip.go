package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Statistics struct {
	success int32
	fail    int32
}
type Bucket struct {
	windowStart int64
	statistics  Statistics
}
type RollingNumber struct {
	timeInMilliseconds      int64
	numberOfBuckets         int64
	bucketSizeInMillseconds int64
}
type ListState struct {
	buckets []*Bucket
	size    int64
	tail    int64
	rollNum *RollingNumber
	mu      sync.RWMutex
}

func NewListState(tail, size, timeInMilliseconds, numberOfBuckets int64) *ListState {
	return &ListState{
		size: size,
		tail: tail,
		rollNum: &RollingNumber{
			timeInMilliseconds:      timeInMilliseconds,
			numberOfBuckets:         numberOfBuckets,
			bucketSizeInMillseconds: timeInMilliseconds / size,
		},
		buckets: make([]*Bucket, size),
	}
}
func (l *ListState) getCurrentBucket() *Bucket {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Unix()
	//最后一个桶在时间窗口内，返回
	currentBucket := l.buckets[int(l.tail)]
	if currentBucket != nil && now < currentBucket.windowStart+l.rollNum.bucketSizeInMillseconds {
		return currentBucket
	}
	//数组为空，新建返回
	if l.tail == 0 && l.buckets[int(l.tail)] == nil {
		bucket := &Bucket{
			windowStart: now,
			statistics: Statistics{
				success: 0,
				fail:    0,
			},
		}
		l.buckets[int(l.tail)] = bucket
		return bucket
	} else {
		var i int64
		for i = 0; i < l.rollNum.numberOfBuckets; i++ {
			lastBucket := l.buckets[int(l.tail)]
			if now < lastBucket.windowStart+l.rollNum.bucketSizeInMillseconds {
				return lastBucket
			} else if now-(lastBucket.windowStart+l.rollNum.bucketSizeInMillseconds) > l.rollNum.timeInMilliseconds {
				l.reset()
				return l.getCurrentBucket()
			} else {
				l.tail++
				bucket := &Bucket{
					statistics:  Statistics{},
					windowStart: lastBucket.windowStart + l.rollNum.bucketSizeInMillseconds,
				}
				if l.tail > l.rollNum.numberOfBuckets {
					l.buckets = l.buckets[1:]
					l.tail--
				}
				l.buckets[int(l.tail)] = bucket
			}
		}
		return l.buckets[int(l.tail)]
	}
}
func (l *ListState) reset() {
	l.buckets = make([]*Bucket, l.size)
	l.tail = 0
}
func (l *ListState) incrSuccess(value int32) {
	bucket := l.getCurrentBucket()
	atomic.AddInt32(&bucket.statistics.success, value)
}
func (l *ListState) incrFail(value int32) {
	bucket := l.getCurrentBucket()
	atomic.AddInt32(&bucket.statistics.fail, value)
}
func (l *ListState) getSum() Statistics {
	var s Statistics
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, bucket := range l.buckets {
		s.success += bucket.statistics.success
		s.fail += bucket.statistics.fail
	}
	return s
}
func main() {
	list := NewListState(0, 5, 50, 10)
	fmt.Println(time.Now().Unix())
	for i := 0; i < 5; i++ {
		list.incrSuccess(int32(i + 1))
		time.Sleep(time.Second * 5)
	}
	for key, bucket := range list.buckets {
		fmt.Printf("bucket[%d]:%+v\n", key, bucket)
	}

}
