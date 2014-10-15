package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
)

var templates = template.Must(template.ParseGlob("templates/*"))
var store = sessions.NewCookieStore([]byte("thisistotallyapasswordimcommittingtogithub"))
var userCount = 0
var userLock = &sync.Mutex{}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler).Methods("GET")
	r.HandleFunc("/", voteHandler).Methods("POST")
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/", r)

	log.Println("Starting server...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
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
	err = session.Save(r, w)
	if err != nil {
		log.Print("Error saving session:", err)
	}
	err = templates.ExecuteTemplate(w, "lively.html", struct{}{})
	if err != nil {
		log.Print("Error executing template:", err)
	}
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "info")
	if err != nil {
		log.Print("Error with getting session:", err)
	}
	session.Values["Vote"] = r.FormValue("Vote")
	session.Save(r, w)
}
