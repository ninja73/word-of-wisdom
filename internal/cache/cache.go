package cache

type Cache interface {
	ContainsOrAdd(val uint64) bool
}
