// Package user provides a simple representation of a user with a password.
package user

type ID = int

type User struct {
	ID       ID
	Name     string
	Password string
}
