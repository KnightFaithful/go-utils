package monotonic

type MonotonicQueue[K comparable] struct {
	list []K
	// 队尾是a，新加入元素是b，b的加入会不会导致a出队列
	popLastCondition func(a, b K) bool
	zero             K
}

func NewMonotonicQueue[K comparable](popLastCondition func(a, b K) bool) *MonotonicQueue[K] {
	return &MonotonicQueue[K]{
		list:             []K{},
		popLastCondition: popLastCondition,
	}
}

func (q *MonotonicQueue[K]) Size() int {
	return len(q.list)
}

func (q *MonotonicQueue[K]) PushLast(val K) {
	for q.Size() > 0 && q.popLastCondition(q.PeekLast(), val) {
		q.PopLast()
	}
	q.list = append(q.list, val)
}

func (q *MonotonicQueue[K]) PopFirst() K {
	if q.Size() == 0 {
		return q.zero
	}
	temp := q.list[0]
	q.list = q.list[1:]
	return temp
}

func (q *MonotonicQueue[K]) PopFirstIfEquals(k K) K {
	if q.Size() == 0 {
		return q.zero
	}
	if q.PeekFirst() == k {
		return q.PopFirst()
	}
	return q.zero
}

func (q *MonotonicQueue[K]) PopLast() K {
	if q.Size() == 0 {
		return q.zero
	}
	temp := q.list[q.Size()-1]
	q.list = q.list[:q.Size()-1]
	return temp
}

func (q *MonotonicQueue[K]) PeekLast() K {
	if len(q.list) == 0 {
		return q.zero
	}
	return q.list[len(q.list)-1]
}

func (q *MonotonicQueue[K]) PeekFirst() K {
	if len(q.list) == 0 {
		return q.zero
	}
	return q.list[0]
}
