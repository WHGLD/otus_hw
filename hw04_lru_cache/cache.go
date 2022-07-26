package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	cachedListItem, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(cachedListItem)
		value := cachedListItem.Value.(cacheItem).value
		return value, ok
	} else {
		return nil, ok
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	cachedListItem, ok := l.items[key]
	if ok {
		cachedListItem.Value = cacheItem{
			key:   key,
			value: value,
		}

		l.queue.MoveToFront(cachedListItem)
	} else {
		currentItemsCount := l.queue.Len()
		if currentItemsCount == l.capacity {
			lastListItem := l.queue.Back()
			l.queue.Remove(lastListItem)
			delete(l.items, lastListItem.Value.(cacheItem).key)
		}

		newListItem := l.queue.PushFront(cacheItem{
			key:   key,
			value: value,
		})
		l.items[key] = newListItem
	}

	return ok
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
