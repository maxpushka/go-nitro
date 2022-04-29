package safesync

import "sync"

// Map wraps sync.Map in order to provide type safety
type Map[T any] struct {
	m sync.Map
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (o *Map[T]) Load(id string) (value T, ok bool) {
	data, ok := o.m.Load(id)

	if !ok {
		var result T
		return result, false
	}

	value = data.(T)

	return value, ok
}

// Store sets the value for a key.
func (o *Map[T]) Store(key string, data T) {
	o.m.Store(key, data)
}

// Delete deletes the value for a key.
func (o *Map[T]) Delete(key string) {
	o.m.Delete(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently, Range may reflect any mapping for that key
// from any point during the Range call.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (o *Map[T]) Range(f func(key string, value T) bool) {
	untypedF := func(key, value interface{}) bool {
		return f(key.(string), value.(T))
	}
	o.m.Range(untypedF)
}
