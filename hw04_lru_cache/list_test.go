package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestPushFront(t *testing.T) {
	l := NewList()

	l.PushFront(10) // [10]
	l.PushFront(20) // [20, 10]
	l.PushFront(30) // [30, 20, 10]
	require.Equal(t, 3, l.Len())

	elems := make([]int, 0, l.Len())
	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}
	require.Equal(t, []int{30, 20, 10}, elems)

	elems2 := make([]int, 0, l.Len())
	for i := l.Back(); i != nil; i = i.Prev {
		elems2 = append(elems2, i.Value.(int))
	}
	require.Equal(t, []int{10, 20, 30}, elems2)
}

func TestPushBack(t *testing.T) {
	l := NewList()

	l.PushBack(10) // [10]
	l.PushBack(20) // [10, 20]
	l.PushBack(30) // [10, 20, 30]
	require.Equal(t, 3, l.Len())

	elems := make([]int, 0, l.Len())
	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}
	require.Equal(t, []int{10, 20, 30}, elems)

	elems2 := make([]int, 0, l.Len())
	for i := l.Back(); i != nil; i = i.Prev {
		elems2 = append(elems2, i.Value.(int))
	}
	require.Equal(t, []int{30, 20, 10}, elems2)
}

func TestRemove(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()
		i := &ListItem{
			Value: 10,
			Next:  &ListItem{},
			Prev:  &ListItem{},
		}

		l.Remove(i) // []
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("only one item", func(t *testing.T) {
		l := NewList()
		i := l.PushFront(10) // [10]

		l.Remove(i) // []
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("remove from front", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)          // [10]
		front := l.PushFront(20) // [20, 10]

		l.Remove(front) // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Nil(t, l.Front().Next)
		require.Nil(t, l.Back().Prev)
	})

	t.Run("remove from back", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)        // [10]
		back := l.PushBack(20) // [10, 20]

		l.Remove(back) // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Back().Value)
		require.Nil(t, l.Front().Next)
		require.Nil(t, l.Back().Prev)
	})

	t.Run("remove from middle", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)           // [10]
		middle := l.PushFront(20) // [20, 10]
		l.PushFront(30)           // [30, 20, 10]

		l.Remove(middle) // [30, 10]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 30, l.Front().Value)
		require.Equal(t, 10, l.Front().Next.Value)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 30, l.Back().Prev.Value)
	})
}

func TestMoveToFront(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()
		i := &ListItem{
			Value: 10,
			Next:  &ListItem{},
			Prev:  &ListItem{},
		}

		l.MoveToFront(i) // []
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("only one item", func(t *testing.T) {
		l := NewList()
		i := l.PushFront(10) // [10]

		l.MoveToFront(i) // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Nil(t, l.Front().Next)
		require.Nil(t, l.Back().Prev)
	})

	t.Run("remove from front", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)          // [10]
		front := l.PushFront(20) // [20, 10]

		l.MoveToFront(front) // [20, 10]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 10, l.Front().Next.Value)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 20, l.Back().Prev.Value)
	})

	t.Run("remove from back", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)        // [10]
		back := l.PushBack(20) // [10, 20]

		l.MoveToFront(back) // [20, 10]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 10, l.Front().Next.Value)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 20, l.Back().Prev.Value)
	})

	t.Run("remove from middle", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)           // [10]
		middle := l.PushFront(20) // [20, 10]
		l.PushFront(30)           // [30, 20, 10]

		l.MoveToFront(middle) // [20, 30, 10]
		require.Equal(t, 3, l.Len())
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 30, l.Front().Next.Value)
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 30, l.Back().Prev.Value)
	})
}
