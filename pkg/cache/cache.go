package cache

//go:generate mockery --dir . --name Cache --filename cache.go --output ./mocks
type Cache interface {
	Get(key interface{}) (interface{}, bool)
	Set(key, value interface{}) bool
}
