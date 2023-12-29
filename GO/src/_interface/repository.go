package _interface

// Repository defines the interface for interacting with user data.
type Repository[T Entity] interface {
	GetByID(id string) (T, error)
	Add(entity T) error
	GetAll() ([]T, error)
	Update(entity T) error
}
