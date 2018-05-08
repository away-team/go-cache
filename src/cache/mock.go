package cache

type mock struct{}

// NewMock returns a Cache that doesn't actually cache anything
func NewMock() Cache {
	return &mock{}
}

func (m *mock) GetAndLoad(key string, loader func() (interface{}, error)) (interface{}, error) {
	return loader()
}
