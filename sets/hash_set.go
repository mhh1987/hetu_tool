package sets

import "sync"

type HashSet[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | string] struct {
	m      map[T]bool
	closed bool
	mu     sync.Mutex
}

func NewHashSet[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | string]() *HashSet[T] {
	return &HashSet[T]{m: make(map[T]bool)}
}

// Add 添加    true 添加成功 false 添加失败
func (set *HashSet[T]) Add(e T) (b bool) {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

// AddAll 添加    true 添加成功 false 添加失败
func (set *HashSet[T]) AddAll(eList ...T) {
	for _, e := range eList {
		set.Add(e)
	}
}

// Remove 删除
func (set *HashSet[T]) Remove(e T) {
	delete(set.m, e)
}

// Clear 清除
func (set *HashSet[T]) Clear() {
	set.m = make(map[T]bool)
}

// Contains 是否包含
func (set *HashSet[T]) Contains(e T) bool {
	return set.m[e]
}

// Len 获取元素数量
func (set *HashSet[T]) Len() int {
	return len(set.m)
}

// KeySets 迭代
func (set *HashSet[T]) KeySets() []T {

	initLen := len(set.m)
	snap := make([]T, initLen)
	actualLen := 0

	for k := range set.m {
		if actualLen < initLen {
			snap[actualLen] = k
		} else {
			snap = append(snap, k)
		}
		actualLen++
	}

	if actualLen < initLen {
		snap = snap[:actualLen]
	}

	return snap
}

func (set *HashSet[T]) ToSlice() []T {
	origin := set.KeySets()
	if len(origin) == 0 {
		return nil
	}
	res := make([]T, 0, len(origin))
	for _, v := range origin {
		res = append(res, v)
	}
	return res
}

func (set *HashSet[T]) Intersect(target *HashSet[T]) *HashSet[T] {

	if target == nil {
		return nil
	}

	if set.Len() == 0 || target.Len() == 0 {
		return nil
	}
	result := NewHashSet[T]()
	for key := range set.m {
		if target.Contains(key) {
			result.Add(key)
		}
	}
	return result
}

func (set *HashSet[T]) Difference(target *HashSet[T]) *HashSet[T] {

	result := NewHashSet[T]()
	if target == nil {
		target = NewHashSet[T]()
	}

	if set.Len() == 0 {
		return result
	}

	for key := range set.m {
		if !target.Contains(key) {
			result.Add(key)
		}
	}
	return result
}

func (set *HashSet[T]) Iter() chan interface{} {

	ch := make(chan interface{})

	go func() {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.closed {
			close(ch)
			return
		}
		for k := range set.m {
			ch <- k
		}
		close(ch)
	}()

	return ch
}

// Equal 判断两个set时候相同
// true 相同 false 不相同
func (set *HashSet[T]) Equal(target *HashSet[T]) bool {

	if target == nil {
		return false
	}
	set.mu.Lock()
	target.mu.Lock()
	defer func() {
		set.mu.Unlock()
		target.mu.Unlock()
	}()

	if set.Len() != target.Len() {
		return false
	}
	for k := range set.m {
		if !target.Contains(k) {
			return false
		}
	}
	return true
}

func (set *HashSet[T]) Close() {
	set.mu.Lock()
	defer set.mu.Unlock()
	set.closed = true
}

func (set *HashSet[T]) Clone() *HashSet[T] {

	result := NewHashSet[T]()
	for t := range set.m {
		result.Add(t)
	}
	return result
}
