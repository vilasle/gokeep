package repository

import "context"

type UserRepository interface {
	Add(context.Context, AddUserDTO) error
	Delete(context.Context, DeleteUserDTO) error
	Get(context.Context, GetUserDTO) (FoundUserDTO, error)
}
