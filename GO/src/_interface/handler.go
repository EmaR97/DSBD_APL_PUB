package _interface

// Handler struct
type Handler[E Entity] struct {
	Repository Repository[E]
}

// NewHandler creates a new Handler with the provided UserService
func NewHandler[E Entity](repository Repository[E]) *Handler[E] {
	return &Handler[E]{
		Repository: repository,
	}
}
