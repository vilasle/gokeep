package model

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

// User - presentation user in storage, keep login by creating user and checking password
type User struct {
	id    int64
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

func (u *User) ID() int64 {
	return u.id
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
	if err := u.r.Add(ctx, u); err != nil {
		return err
	}
	return nil

}

// update - update user in storage
func (u *User) update(ctx context.Context) error {
	if err := u.r.Update(ctx, u); err != nil {
		return err
	}
	return nil
}

// delete - delete user from storage
func (u *User) delete(ctx context.Context) error {
	if err := u.r.Delete(ctx, u); err != nil {
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
	return r.Find(ctx, login)
}

func getUserByID(ctx context.Context, id int64, r UserRepository) (*User, error) {
	return r.Get(ctx, id)
}
