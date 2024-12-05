package main

import (
	"sync"
	"sync/atomic"
)

type Counter interface {
	Inc()
	Load() int64
}

// 普通版，非并发安全
type CommonCounter struct {
	counter int64
}

func (c CommonCounter) Inc() {
	c.counter++
}
func (c CommonCounter) Load() int64 {
	return c.counter
}

// 互斥锁版，并发安全
type MutexCounter struct {
	counter int64
	lock    sync.Mutex
}

func (m *MutexCounter) Inc() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.counter++
}
func (m *MutexCounter) Load() int64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.counter
}

// 原子操作版，并发安全且比互斥锁效率高
type AtomicCounter struct {
	counter int64
}

func (a *AtomicCounter) Inc() {
	atomic.AddInt64(&a.counter, 1)
}
func (a *AtomicCounter) Load() int64 {
	return atomic.LoadInt64(&a.counter)
}
