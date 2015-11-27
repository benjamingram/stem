package hosts

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProfile struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

var (
	loginTemplate = template.Must(template.ParseFiles("hosts/launcher-login.html"))
	conf          *oauth2.Config
)

// LoadSecretsFromFile loads secret from file
func LoadSecretsFromFile(filename string) {
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}

	conf, err = google.ConfigFromJSON(configFile, scopes...)

	if err != nil {
		log.Fatal(err)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	loginTemplate.Execute(w, nil)
}

func loginOauthHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login-oauth" {
		http.NotFound(w, r)
		return
	}

	url := conf.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func loginOauthResponseHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)

	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(oauth2.NoContext, tok)

	response, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil || response.StatusCode != 200 {
		log.Fatal(err)
		return
	}

	defer response.Body.Close()

	var profile GoogleProfile

	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&profile); err != nil {
		log.Fatal(err)
		return
	}

	log.Println(profile)
}
