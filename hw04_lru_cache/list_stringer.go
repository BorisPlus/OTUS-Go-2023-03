package hw04lrucache

import (
	"fmt"
)

// String - наглядное представление значения элемента двусвязного списка.
//
// Например,
//
//	-------------------             -------------------
//	Item: 0xc00002e400              Item: 0xc00002e400
//	-------------------             -------------------
//	Data: 30                или     Data: 30
//	Prev: 0xc00002e3c0              Prev: 0x0
//	Next: 0xc00002e440              Next: 0x0
//	-------------------             -------------------
func (listItem *ListItem) String() string {
	template := `
-------------------
Item: %p
-------------------
Data: %v
Prev: %p
Next: %p
-------------------`
	return fmt.Sprintf(template, listItem, listItem.Data, listItem.Prev, listItem.Next)
}

// String - наглядное представление всего двусвязного списка.
//
// Например,
//
// - пустой список:
//
//	(nil:0x0)
//	    |
//	    V
//	(nil:0x0)
//
// - список из двух элементов:
//
//	    (nil:0x0)
//	        |
//	        V
//	-------------------
//	Item: 0xc00002e3a0 <--------┐
//	-------------------         |
//	Data: 2                     |
//	Prev: 0x0                   |
//	Next: 0xc00002e380	>>>-----|---┐
//	-------------------         |   | ссылается на...
//	        |                   |   |
//	        V                   |   |
//	-------------------         |   |
//	Item: 0xc00002e380  <-----------┘
//	-------------------         |
//	Data: 1                     | ссылается на...
//	Prev: 0xc00002e3a0	>>>-----┘
//	Next: 0x0
//	-------------------
//	        |
//	        V
//	    (nil:0x0)
func (list *List) String() string {
	result := ""
	Nill := `
    (nil:0x0)`
	delimiter := `
	|
	V`
	result += Nill
	result += delimiter
	for i := list.Front(); i != nil; i = i.Next {
		result += i.String()
		result += delimiter
	}
	result += Nill
	return result
}
