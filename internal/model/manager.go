package model

type ModelType = int8

const (
	ModelTypeUser ModelType = iota + 1
)

// ModelManager is a manager for all models.
type ModelManager struct {
	repository RepositoryCollector
	User       *userManager
}

// NewModelManager creates a new model manager
func NewModelManager(repository RepositoryCollector) *ModelManager {
	return &ModelManager{
		repository: repository,
		User: &userManager{
			repository: repository.User(),
		},
	}
}

// userManager is a manager for users
type userManager struct {
	repository UserRepository
}

// New creates a new user with login and hashed password
func (c *userManager) New(login, password string) *User {
	return newUser(login, password, c.repository)
}

// FindByLogin  finds a user by login on repository, return ErrUserNotFound if not found
func (c *userManager) FindByLogin(login string) (*User, error) {
	//TODO add logger
	return findUserByLogin(login, c.repository)
}
