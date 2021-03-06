package main

/*
In Config:

    AuthCASUrl        // CAS server
    AuthJWTTimeOut    // int: Hours for jwt timeout
    AuthJWTPassword   // JWT secret password
    AuthJWTCallback   // client url callback to validate and register jwt
    AuthValidLogins   // array of valid cas users


In gin server:

    casHandler := initCas("/admin", config)
    // add /auth/login /auth/logout,  allow cas login to set jwt token
    auth := r.Group("auth")
    {
        auth.GET("/login", gin.WrapH(casHandler))
        auth.GET("/logout", gin.WrapH(casHandler))
    }

*/

import (
	"fmt"
	"log"
	"net/http"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"

	//	"gopkg.in/cas.v2"
	"github.com/go-cas/cas"
	"net/url"
)

type myCasHandler struct{ Config Config }

const error500 = `<!DOCTYPE html>
<html>
  <head>
    <title>Error 500</title>
  </head>
  <body>
    <h1>Error 500</h1>
    <p>%v</p>
  </body>
</html>`

func (h *myCasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := h.Config.Verbose
	ip := r.Header.Get("X-Forwarded-For")
	//fmt.Println("CAS Login")
	if !cas.IsAuthenticated(r) {
		if v == true {
			fmt.Println("CAS Login redirect")
		}
		cas.RedirectToLogin(w, r)
		return
	}

	if r.URL.Path == "/logout" {
		log.Printf("%s: CAS Logout\n", ip)
		cas.RedirectToLogout(w, r)
		return
	}

	username := cas.Username(r)
	log.Printf("%s [%s]: CAS Login\n", ip, username)

	if !contains(h.Config.AuthValidLogins, username) {
		log.Printf("%s [%s]: CAS Not admin access\n", ip, username)
		cas.RedirectToLogout(w, r)
		return
	}

	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	// Set some claims
	token.Claims = jwt_lib.MapClaims{
		"id":   username,
		"role": "admin", //TODO add other roles
		"exp":  time.Now().Add(time.Hour * time.Duration(h.Config.AuthJWTTimeOut)).Unix(),
	}

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(h.Config.AuthJWTPassword))
	if err != nil { // mainly timeout
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s: CAS jwt err\n", ip)
		fmt.Fprintf(w, error500, err)
		return
	}

	// Redirect to callback with token and whithout cache
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
	w.WriteHeader(http.StatusOK)
	redirect := fmt.Sprintf("%s?%s", h.Config.AuthJWTCallback, tokenString)
	http.Redirect(w, r, redirect, 301)

}

func setCasHandler(config Config) http.Handler {
	mh := http.NewServeMux()
	CasHandler := &myCasHandler{}
	CasHandler.Config = config
	mh.Handle("/", CasHandler)
	u, _ := url.Parse(config.AuthCASUrl)
	client := cas.NewClient(&cas.Options{
		URL:         u,
		SendService: true,
	})

	return client.Handle(mh)
}
