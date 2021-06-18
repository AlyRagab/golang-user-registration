package models

import (
	"errors"
	"regexp"
	"strings"

	"github.com/AlyRagab/golang-user-registration/modules/rand"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Create method to create a user in database
func (uv *userValidator) Create(user *User) error {
	// Validating and Normalizing then creating the hash
	if err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.hmacRememberToken,
		//uv.emialIsExisted,
		uv.normalizeEmail); err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Update will hash a remember token if it is provided.
func (uv *userValidator) Update(user *User) error {
	// Validating and Normalizing then updating the hash
	if err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.hmacRememberToken,
		//uv.emialIsExisted,
		uv.normalizeEmail); err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

// ByEmail will normalize the Email address before calling
// ByEmail on the UserDB field
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	// Validating and Normalizing the Email
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// Delete method to delete a user in database with the provided ID only
func (uv *userValidator) Delete(id uint) error {
	user := User{
		Model: gorm.Model{
			ID: id,
		},
	}
	// Validating if the user.ID <= 0
	if err := runUserValFuncs(&user, uv.isGreaterThan(0)); err != nil {
		return err
	}
	user = User{Model: gorm.Model{ID: id}}
	return uv.UserDB.Delete(user.ID)
}

// ByRemember method will hash the RememberToken
// And call the other ByRemember method for making UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	// Validating and Normalizing then creating the hash
	if err := runUserValFuncs(&user, uv.hmacRememberToken); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	return nil
}

// Normalizing Email :
// "Triming space"  and "Lower Case" and "Require Email" and regexp
func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)

	if user.Email == "" {
		return errors.New("Email Address is required")
	}

	// Validate the Email Format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("Email Address is not valid")
	}
	return nil
}

func (uv *userValidator) emialIsExisted(user *User) error {
	// Validating if the Email is already existed
	existed, err := uv.ByEmail(user.Email)
	if err != nil {
		return err
	}
	if user.ID != existed.ID {
		return errors.New("Email Address is already taken")
	}
	return nil
}

// Validating and Normalizing
func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// bcryptPassword will hash the user passsword with
// salt and pepper and bcrypt the password
func (uv *userValidator) bcryptPassword(user *User) error {
	pwPepper := []byte(user.Password + UserPwPepper) // Salt + Pepper
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pwPepper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes) // store the hashedBytes in the struct
	user.Password = ""                      // Don't store Password

	// Look at the token if it is empty then create it.
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	return nil
}

// Validating the password length should be atleast 8 charachters
func (uv *userValidator) passwordLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return errors.New("Passowrd should be atleast 8 char")
	}
	return nil
}

// Validating thet Password is Required
func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return errors.New("Password is required")
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return errors.New("Hashing Password is required")
	}
	return nil
}

// Validating the RememberToken and hashing
func (uv *userValidator) hmacRememberToken(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hashing(user.Remember)
	return nil
}

// Validating if the user.ID is greater than 0
func (uv *userValidator) isGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrInvalidID
		}
		return nil
	})
}
