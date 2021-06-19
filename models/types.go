package models

import (
	"errors"

	"github.com/AlyRagab/golang-user-registration/modules/hash"
	"github.com/jinzhu/gorm"
)

var (
	// ErrNotFound is returned if the resource is not found in the database.
	ErrNotFound = errors.New("Models: Resource Not Found")
	// ErrInvalidID is used when we pass an ID to Delete Method to delete a user from DB
	ErrInvalidID = errors.New("Models: ID must be Valid ID")
	// UserPwPepper Adding the Pepper value
	UserPwPepper = "secret-random-string"
	// ErrInvalidPassword to return Invalid Password
	ErrInvalidPassword = errors.New("Models: Invalid Password")
	// HmacSecret for creating the HMAC
	HmacSecret        = "secret-hmac-key"
	_          UserDB = &userGorm{}
)

// We will implement 3 layers for interacting with DB
// 1. Single CRUD operations Layer
// 2. Authentication Layer
// 3. Normalization and Validation Layer

// UserDB interface handling all User Operations in the database
// This is the Database Layer for single user queries
type UserDB interface {
	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Close DB connection
	Close() error

	// Migration Helpers
	DBDestructiveReset()

	// Handling Database Communication
	Ping() error
}

// UserService interface is a ser of methods used to work with the user model
// Database Auth Layer
type UserService interface {
	// Authenticate will verify the provided email and password
	// If corresponds then the user to that email will be returned
	// else it will return either :
	// ErrNotFound , ErrInvalidPassword or error something worng
	Authenticate(email, password string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
}

// Validation in getting anything from DB
type userValidator struct {
	UserDB
	hmac hash.HMAC
}

type userGorm struct {
	db *gorm.DB
}

// User Model struct
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"` // Not to store in Database
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

type userValFunc func(*User) error
