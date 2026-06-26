package storage

import "github.com/aman-void/go-http-server/internal/types"

type Storage interface {
	CreateUser(name string, email string, age int) (int64, error)
	GetUserById(id int64) (types.User, error)
	GetUsers() ([]types.User, error)
}
