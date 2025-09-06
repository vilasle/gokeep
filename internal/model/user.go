package model

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type UserAccess interface {
	ID() int
	Login() string
	Hash() string
	PasswordIsValid(password string) bool
}

// User - presentation user in storage, keep login by creating user and checking password
type User struct {
	id    int
	login string
	//hashed password
	password string
	r        UserRepository
}

// newUser - create new user and create password hash
func newUser(login, password string, r UserRepository) *User {
	return &User{
		login:    login,
		password: hash256(password),
		r:        r,
	}
}

// PasswordIsValid - check matching passwords
func (u *User) PasswordIsValid(password string) bool {
	return u.password == hash256(password)
}

//ID - return user id
func (u *User) ID() int {
	return u.id
}

//Hash - return user password hash
func (u *User) Hash() string {
	return u.password
}

//Login - return user login
func (u *User) Login() string {
	return u.login
}

// Save - save user to storage
func (u *User) Save(ctx context.Context) error {
	if u.isExists() {
		return u.update(ctx)
	}
	return u.add(ctx)
}

// Delete - delete user from storage
func (u *User) Delete(ctx context.Context) error {
	if u.isExists() {
		return u.delete(ctx)
	}
	return ErrUserNotFound
}

// isExists - return false if user does not exists in storage
func (u *User) isExists() bool {
	return u.id > 0
}

// add - add new user to storage
func (u *User) add(ctx context.Context) error {
	if id, err := u.r.Add(ctx, UserAdd{Login: u.login, Password: u.password}); err != nil {
		return err
	} else {
		u.id = id
	}
	return nil

}

// update - update user in storage
func (u *User) update(ctx context.Context) error {
	if err := u.r.Update(ctx, UserUpdate{ID: u.id, Login: u.login, Password: u.password}); err != nil {
		return err
	}
	return nil
}

// delete - delete user from storage
func (u *User) delete(ctx context.Context) error {
	if err := u.r.Delete(ctx, u.id); err != nil {
		return err
	}
	return nil

}

// hash - hash string with sha256
func hash256(v string) string {
	hasher := sha256.New()
	//ignore error because sha256.Digest.Write() does not return error
	//and has function signature ([]byte)(int, error) for implementation io.Writer
	_, _ = hasher.Write([]byte(v))

	return hex.EncodeToString(hasher.Sum(nil))
}

func findUserByLogin(ctx context.Context, login string, r UserRepository) (*User, error) {
	if info, err := r.Find(ctx, login); err == nil {
		return &User{
			id:       info.ID,
			login:    info.Login,
			password: info.Password,
			r:        r,
		}, nil
	} else {
		return nil, err
	}
}

func getUserByID(ctx context.Context, id int, r UserRepository) (*User, error) {
	if info, err := r.Get(ctx, id); err == nil {
		return &User{
			id:       info.ID,
			login:    info.Login,
			password: info.Password,
			r:        r,
		}, nil
	} else {
		return nil, err
	}
}
