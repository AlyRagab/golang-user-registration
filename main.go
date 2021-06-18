package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AlyRagab/golang-user-registration/controllers"
	"github.com/AlyRagab/golang-user-registration/models"
	"github.com/hellofresh/health-go/v4"
	healthPg "github.com/hellofresh/health-go/v4/checks/postgres"

	"github.com/gorilla/mux"
)

// Error Handling
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// HTTP 404 NotFound
func notfound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 Page not found")
}

var psqlInfo string

func main() {
	var (
		host     = os.Getenv("DB_HOST")
		port     = 5432
		dbuser   = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, dbuser, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	//us.DBDestructiveReset()

	staticC := controllers.NewStatic() // Parsing static templates
	usersC := controllers.NewUsers(us) // Handling User Controller

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(notfound)

	// healthz endpoint
	h, _ := health.New()
	h.Register(health.Config{
		Name:      "postgres-check",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: healthPg.New(healthPg.Config{
			DSN: psqlInfo,
		}),
	})
	r.Handle("/healthz", h.Handler())
	fmt.Println("Starting Server on 0.0.0.0:8080")
	http.ListenAndServe(":8080", r)
}
