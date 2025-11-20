package storage

import "gopher/internal/model"

type UserRepository interface {
	Create(*model.User) error
	FindById(string) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
