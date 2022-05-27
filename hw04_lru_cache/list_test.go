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

	t.Run("push", func(t *testing.T) {
		l := NewList()
		middleItem := l.PushFront("middle")

		l.PushFront("front")
		require.Equal(t, "front", l.Front().Value)
		require.Nil(t, l.Front().Prev)
		require.Equal(t, middleItem, l.Front().Next)

		l.PushBack("back")
		require.Equal(t, "back", l.Back().Value)
		require.Nil(t, l.Back().Next)
		require.Equal(t, middleItem, l.Back().Prev)
	})

	t.Run("remove", func(t *testing.T) {
		l := NewList()
		first := l.PushFront(1)
		l.PushBack(2)
		middle := l.PushBack(3)
		last := l.PushBack(4)

		l.Remove(first) // [2, 3, 4]
		require.Equal(t, 2, l.Front().Value)

		l.Remove(middle) // [2, 4]
		require.Equal(t, 2, l.Front().Value)
		require.Equal(t, 4, l.Back().Value)

		l.Remove(last)
		require.Equal(t, 2, l.Back().Value)
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
