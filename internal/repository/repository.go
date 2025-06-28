package repository

import "context"

type UserRepository interface {
	Add(context.Context, AddUserDTO) error
	Get(context.Context, GetUserDTO) (FoundUserDTO, error)
}
