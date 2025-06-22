package repository

import "time"

type AddUserDTO struct {
	Login    string
	Password string
}

type GetUserDTO struct {
	Login string
}

type FoundUserDTO struct {
	ID        int
	Login     string
	Password  string
	CreatedAt time.Time
}

type UpdateUserDTO struct {
	ID       int
	Login    string
	Password string
}

type DeleteUserDTO struct {
	ID int
}
