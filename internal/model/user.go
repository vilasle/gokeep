package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// User - presentation user in storage, keep login by creating user and checking password
type User struct {
	ID    int64
	Login string
	//hashed password
	Password  string
	CreatedAt time.Time
	r         UserRepository
}

// newUser - create new user and create password hash
func newUser(login, password string, r UserRepository) *User {
	return &User{
		Login:    login,
		Password: hash256(password),
		r:        r,
	}
}

// PasswordIsValid - check matching passwords
func (u *User) PasswordIsValid(password string) bool {
	return u.Password == hash256(password)
}

// IsExists - return false if user does not exists in storage
func (u *User) IsExists() bool {
	return u.ID > 0
}

// Save - save user to storage
func (u *User) Save() error {
	if u.IsExists() {
		return u.update()
	}
	return u.add()
}

// Delete - delete user from storage
func (u *User) Delete() error {
	if !u.IsExists() {
		return ErrUserNotExists
	}
	return u.delete()
}

// add - add new user to storage
func (u *User) add() error {
	if err := u.r.Add(u); err != nil {
		return err
	}
	return nil

}

// update - update user in storage
func (u *User) update() error {
	if err := u.r.Update(u); err != nil {
		return err
	}
	return nil
}

// delete - delete user from storage
func (u *User) delete() error {
	if err := u.r.Delete(u); err != nil {
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

func findUserByLogin(login string, r UserRepository) (*User, error) {
	return r.Get(login)
}
