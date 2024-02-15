package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

func main() {

	scope := "user-read-recently-played"
	clientID := ""
	clientSecret := ""
	redirectURI := "http://localhost:8080/callback"

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     spotify.Endpoint,
		Scopes:       []string{scope},
		RedirectURL:  redirectURI,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test\n")
	})

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, config.AuthCodeURL(""), http.StatusFound)
	})

	mux.HandleFunc("GET /callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := config.Exchange(r.Context(), code)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s\n", token.RefreshToken)
	})

	log.Println("starting server...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
