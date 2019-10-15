package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func newRouter(ds dataStore) *mux.Router {

	r := mux.NewRouter()

	r.Handle("/messages", getMessages(ds)).Methods("GET")
	r.Handle("/messages", newMessage(ds)).Methods("POST")
	r.Handle("/messages/{id:[0-9]+}", updateMessage(ds)).Methods("PUT")
	r.Handle("/messages/{id:[0-9]+}", getMessage(ds)).Methods("GET")
	r.Handle("/messages/{id:[0-9]+}", deletMessage(ds)).Methods("DELETE")
	r.Handle("/auth/{provider}/callback", providerCallback())
	r.Handle("/logout/{provider}", logout())
	r.Handle("/auth/{provider}", login())
	r.Handle("/login", index())

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	return r

}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("login called")
		t, _ := template.New("foo").Parse(indexTemplate)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t.Execute(w, providerIndex)
	})
}

func login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
			t, _ := template.New("foo").Parse(userTemplate)
			t.Execute(w, gothUser)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	})
}

func logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}

func providerCallback() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(w, user)
	})
}

func getMessages(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err, msg := db.GetMessages()
		log.Println(msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(msg); err != nil {
				panic(err)
			}
		}
	})

}

func getMessage(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			_, msg := db.GetMessage(id)
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(msg); err != nil {
				log.Println("falid to parse json", err)
			}
		}
	})
}

func newMessage(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		bodyString := string(body)
		log.Printf("new message: %s", bodyString)
		if err := db.NewMessage(bodyString); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
}

func updateMessage(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		body, _ := ioutil.ReadAll(r.Body)
		bodyString := string(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			db.UpdateMessage(id, bodyString)
			w.WriteHeader(http.StatusOK)
		}
	})
}

func deletMessage(db dataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			db.DeleteMessage(id)
			w.WriteHeader(http.StatusOK)
		}
	})
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <r><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></r>
{{end}}`

var userTemplate = `
<r><a href="/logout/{{.Provider}}">logout</a></r>
<r>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</r>
<r>Email: {{.Email}}</r>
<r>NickName: {{.NickName}}</r>
<r>Location: {{.Location}}</r>
<r>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></r>
<r>Description: {{.Description}}</r>
<r>UserID: {{.UserID}}</r>
<r>AccessToken: {{.AccessToken}}</r>
<r>ExpiresAt: {{.ExpiresAt}}</r>
<r>RefreshToken: {{.RefreshToken}}</r>
`
