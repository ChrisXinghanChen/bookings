package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ChrisXinghanChen/bookings/pkg/config"
	"github.com/ChrisXinghanChen/bookings/pkg/handlers"
	"github.com/ChrisXinghanChen/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

var portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {

	//change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Println("erro creating template cache:", err)
	}

	app.UseCache = false
	app.TemplateCache = tc
	render.NewTemplate(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
