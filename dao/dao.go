package dao

type Model[T any] struct {
	Result T
	Error  error
}

type Datastore[T Model[T]] interface {
	Create(value T) error
	Update(value T) error
	Delete(value T) error
	Find(dest interface{}, conds ...interface{}) (T, error)
	FindAll(dest interface{}, conds ...interface{}) ([]T, error)
}
