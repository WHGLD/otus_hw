package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	First  *ListItem
	Last   *ListItem
	Length int
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.First
}

func (l *list) Back() *ListItem {
	return l.Last
}

// PushFront Добавляем в начало
func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{
		v,
		nil,
		nil,
	}

	if l.Len() == 0 {
		l.First = &item
		l.Last = l.First
	} else {
		item.Next = l.First
		l.First.Prev = &item
		l.First = &item
	}

	l.Length++

	return &item
}

// PushBack Добавляем в конец
func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{
		v,
		nil,
		nil,
	}

	if l.Len() == 0 {
		l.First = &item
		l.Last = l.First
	} else {
		item.Prev = l.Last
		l.Last.Next = &item
		l.Last = &item
	}

	l.Length++

	return &item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		// надо переписать next у prev
		if i.Next != nil {
			i.Prev.Next = i.Next
		} else {
			i.Prev.Next = nil
		}
	} else {
		if l.Len() > 1 {
			l.First = i.Next
		} else {
			l.First = nil
		}
	}

	if i.Next != nil {
		// надо переписать prev у next
		if i.Prev != nil {
			i.Next.Prev = i.Prev
		} else {
			i.Next.Prev = nil
		}
	} else {
		if l.Len() > 1 {
			l.Last = i.Prev
		} else {
			l.Last = nil
		}
	}

	i.Prev = nil
	i.Next = nil

	l.Length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.First != i && l.Len() > 1 {
		l.Remove(i)
		// сохраним указатель на бывший первый элемент
		firstTemp := l.First
		// обновим ссылки у перемещаемого элемента
		i.Prev = nil
		i.Next = firstTemp
		// актуализируем ссылки бывшего первого элемента
		firstTemp.Prev = i
		// переназначим первый элемент
		l.First = i
	}

}

func NewList() List {
	return new(list)
}
