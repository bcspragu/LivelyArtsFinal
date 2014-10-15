package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"sync"
)

var templates = template.Must(template.ParseGlob("templates/*"))
var store = sessions.NewCookieStore([]byte("thisistotallyapasswordimcommittingtogithub"))
var userCount = 0
var userLock = &sync.Mutex{}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler).Methods("POST")
	r.HandleFunc("/", voteHandler).Methods("GET")
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/", r)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "info")
	if err != nil {
		log.Print("Error with getting session:", err)
	}
	if session.IsNew {
		userLock.Lock()
		userCount++
		session.Values["ID"] = userCount
		session.Values["Vote"] = nil
		userLock.Unlock()
	}
	session.Save(r, w)
	templates.ExecuteTemplate(w, "lively.html", struct{}{})
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "info")
	if err != nil {
		log.Print("Error with getting session:", err)
	}
	session.Values["Vote"] = r.FormValue("Vote")
	session.Save(r, w)
}
