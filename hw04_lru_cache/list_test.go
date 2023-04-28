package hw04lrucache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// go test -v list.go list_string.go cache.go
func TestList(t *testing.T) {

	t.Run("zero-value list-item test", func(t *testing.T) {
		zeroValueItem := ListItem{}
		require.Nil(t, zeroValueItem.Data)
		require.Nil(t, zeroValueItem.Prev)
		require.Nil(t, zeroValueItem.Next)
		fmt.Println("\n[zero-value] is:\n", &zeroValueItem)
	})

	t.Run("list-items referencies test", func(t *testing.T) {

		fmt.Println("\n[1] <--> [2] <--> [3]")
		first := &ListItem{
			Data: 1,
			Prev: nil,
			Next: nil,
		}
		fmt.Println("\n[1] is:\n", first)

		second := &ListItem{
			Data: 2,
			Prev: first,
			Next: nil,
		}
		first.Next = second
		fmt.Println("\n[1] become:\n", second)
		fmt.Println("\n[2] is:\n", second)

		third := ListItem{
			Data: 3,
			Prev: second,
			Next: nil,
		}

		second.Next = &third
		fmt.Println("\n[2] become:\n", second)
		fmt.Println("\n[3] is:\n", &third)

		fmt.Println("first.Next.Next.Next is nil")
		require.Nil(t, first.Next.Next.Next)
		fmt.Println("first.Next.Next is third")
		require.Equal(t, &third, first.Next.Next)
		fmt.Println("third.Prev.Prev.Prev is nil")
		require.Nil(t, third.Prev.Prev.Prev)
		fmt.Println("third.Prev.Prev is first")
		require.Equal(t, first, third.Prev.Prev)
	})

	t.Run("empty list test", func(t *testing.T) {
		list := NewList()

		require.Equal(t, 0, list.Len())
		require.Nil(t, list.Front())
		require.Nil(t, list.Back())
	})

	t.Run("little list test #1", func(t *testing.T) {
		list := NewList()

		fmt.Println("\nList was:\n", list)

		item := list.PushFront(1) // [1]
		fmt.Println("\nItem was:\n", item)

		list.MoveToFront(item)
		fmt.Println("\nItem become:\n", item)

		fmt.Println("\nList become:\n", list)

		fmt.Println("\nBack become:\n", list.Back())
		fmt.Println("\nFront become:\n", list.Front())
	})

	t.Run("little list test #2", func(t *testing.T) {
		list := NewList()
		fmt.Println("\nList was:\n", list)

		item := list.PushFront(1) // [1]
		fmt.Println("\nItem [1] become:\n", item)

		back := list.Back()
		fmt.Println("\nItem [back] become:\n", back)
		front := list.Front()
		fmt.Println("\nItem [front] become:\n", front)

		item = list.PushFront(2) // [2]
		fmt.Println("\nItem [2] become:\n", item)

		back = list.Back()
		fmt.Println("\nItem [back] become:\n", back)
		front = list.Front()
		fmt.Println("\nItem [front] become:\n", front)

		fmt.Println("\nList become:\n", list)

		list.Remove(back)
		fmt.Println("\nItem [back] removed:\n", back)

		fmt.Println("\nList become:\n", list)

		require.Equal(t, list.Front(), list.Back())
		require.Nil(t, list.Front().Prev)
		require.Nil(t, list.Front().Next)

		back = list.Back()
		fmt.Println("\nItem [back] become 2 step:\n", back)
		list.Remove(back)
		fmt.Println("\nItem [back] removed  2 step:\n", back)

		fmt.Println("\nList become:\n", list)

		require.Equal(t, list.Front(), list.Back())
		require.Nil(t, list.Front())
		require.Nil(t, list.Back())
	})

	t.Run("complex processing test", func(t *testing.T) {
		list := NewList()

		list.PushFront(10) // [10]
		fmt.Println("\n[10] become:\n", list)
		list.PushBack(20) // [10, 20]
		fmt.Println("\n[10, 20] become:\n", list)
		list.PushBack(30) // [10, 20, 30]
		fmt.Println("\n[10, 20, 30] become:\n", list)
		require.Equal(t, 3, list.Len())

		middle := list.Front().Next // 20
		require.Equal(t, middle.Data, 20)
		fmt.Printf("middle.Value is %v. OK\n", middle.Data)
		list.Remove(middle) // [10, 30]
		require.Equal(t, 2, list.Len())
		fmt.Println("[10, 30] become:\n", list)

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				list.PushFront(v)
			} else {
				list.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]
		fmt.Println("[10, 30] mix [40, 50, 60, 70, 80] with mod(2, idx)")

		require.Equal(t, 7, list.Len())
		fmt.Printf("list.Len() is %v. OK\n", list.Len())
		require.Equal(t, 80, list.Front().Data)
		fmt.Printf("list.Front().Value is %v. OK\n", list.Front().Data)
		require.Equal(t, 70, list.Back().Data)
		fmt.Printf("list.Back().Value is %v. OK\n", list.Back().Data)

		// Check for [80, 60, 40, 10, 30, 50, 70]
		var elems []int
		// - forward stroke
		elems = make([]int, 0, list.Len())
		for i := list.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Data.(int))
		}
		require.Equal(t, []int{80, 60, 40, 10, 30, 50, 70}, elems)
		fmt.Println("Forward stroke check for [80, 60, 40, 10, 30, 50, 70]. OK")
		// - reverse stroke
		elems = make([]int, 0, list.Len())
		for i := list.Back(); i != nil; i = i.Prev {
			elems = append([]int{i.Data.(int)}, elems...)
		}
		require.Equal(t, []int{80, 60, 40, 10, 30, 50, 70}, elems)
		fmt.Println("Reverse stroke check for [80, 60, 40, 10, 30, 50, 70]. OK")
		// - move front to front
		list.MoveToFront(list.Front())
		elems = make([]int, 0, list.Len())
		for i := list.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Data.(int))
		}
		require.Equal(t, []int{80, 60, 40, 10, 30, 50, 70}, elems)
		fmt.Println("Move front to front check for [80, 60, 40, 10, 30, 50, 70]. OK")
		// - remove last and put it to last
		last := list.Back()
		list.Remove(last)
		list.PushBack(last.Data)
		elems = make([]int, 0, list.Len())
		for i := list.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Data.(int))
		}
		require.Equal(t, []int{80, 60, 40, 10, 30, 50, 70}, elems)
		fmt.Println("Remove and PushBack last - check for [80, 60, 40, 10, 30, 50, 70]. OK")
		// - check for nil-refs of first and last
		require.Nil(t, list.Front().Prev)
		require.Nil(t, list.Back().Next)
		fmt.Println("Check for list.Front().Prev and list.Back().Next is nils. OK")

		list.MoveToFront(list.Back()) // [70, 80, 60, 40, 10, 30, 50]
		// Check for [70, 80, 60, 40, 10, 30, 50]
		elems = make([]int, 0, list.Len())
		for i := list.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Data.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
		fmt.Println("Move back to front check for [80, 60, 40, 10, 30, 50, 70]. OK")

		fmt.Printf("list become\n%s\n", list)
	})
}
