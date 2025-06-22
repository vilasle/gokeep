package model

type RepositoryCollector interface {
	User() UserRepository
}

type UserRepository interface {
	Add(*User) error
	Update(*User) error
	Delete(*User) error
	Get(string) (*User, error)
}
