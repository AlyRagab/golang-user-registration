package controllers

import (
	"fmt"
	"net/http"

	"github.com/AlyRagab/golang-user-registration/models"
	"github.com/AlyRagab/golang-user-registration/modules/rand"
	"github.com/AlyRagab/golang-user-registration/views"
)

// NewUsers for Parsing new user view/template in signup page
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

// New method to Render and Execute the new user view/template
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create method for creating new user " POST /signup" endpoint
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := u.signIn(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect to other page after the login
	http.Redirect(w, r, "/login", http.StatusFound)
}

// Login is a POST /login endpoint
// Verify the provided email and password are correct then login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid Email Address")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid Password")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect to other page after the login
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// signIn helps in setting the cookie "email" to the end user
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	// Making sure to Remember the token
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true, // means that it is not accessible to scripts "to protect against XSS"
	}
	http.SetCookie(w, &cookie)
	return nil
}

// CookieTest is used to display cookies set on the current user :
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "The user is %+v", user)
	//fmt.Fprintln(w, "Email is :", cookie.Value)
}
