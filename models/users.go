package models

import (
	"errors"

	"github.com/AlyRagab/golang-user-registration/modules/hash"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewUserService func for creating connection to database
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(HmacSecret)
	uv := &userValidator{
		UserDB: ug,
		hmac:   hmac,
	}
	return &userService{
		UserDB: uv,
	}, nil
}

// Implementing and returning the userGorm type
func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db: db,
	}, nil
}

func (ug *userGorm) Ping() error {
	if err := ug.db.DB().Ping(); err != nil {
		ug.db.DB().Close()
		return errors.New("Connection to DB is not available")
	}
	return nil
}

// ByID method to get a user by ID
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id).First(&user)
	err := first(db, &user)
	return &user, err
}

// ByEmail method to get a user by Email
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email).First(&user)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with the given remember token
// and returns that user. This mdethos expects the remember
// token to be already hashed.
// Errors handeled as the same done by the ByEmail.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// will query the gorm.DB and get the first item from db and place it into
// dst , if nothing is found , it will return error.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// CloseDB to be used as `defer us.db.Close()`
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DBDestructiveReset drops the user table and rebuilds it again - DEV only :)
func (ug *userGorm) DBDestructiveReset() {
	ug.db.DropTableIfExists(&User{})
	ug.db.AutoMigrate(&User{})
}

// Create method to create a user in database
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update method to update a user in database
func (ug *userGorm) Update(user *User) error {
	return ug.db.Model(&user).Where("name = ?", &user.Name).Update(&user).Error
}

// Delete method to delete a user in database with the provided ID only
func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Where("id = ?", id).Delete(&user).Error
}

// Authenticate Method is used for Authenticate and Validate login
func (us *userService) Authenticate(email, password string) (*User, error) {
	// Vlidate if the user is existed in the database or no
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	// Compare the login based in the Hash value
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+UserPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		case nil:
			return nil, err
		default:
			return nil, err
		}
	}
	return foundUser, nil
}
