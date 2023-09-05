package parser

func NewLoopedList[T any](size int) LoopedList[T] {
	return LoopedList[T]{
		data:  make([]T, size),
		start: 0,
	}
}

type LoopedList[T any] struct {
	data  []T
	start int
}

func (l LoopedList[T]) Size() int {
	return len(l.data)
}

func (l LoopedList[T]) Get(i int) T {
	return l.data[(l.start+i)%l.Size()]
}

func (l *LoopedList[T]) Push(val T) {
	l.data[l.start] = val
	l.start = (l.start + 1) % l.Size()
}
